package repository

import (
	"AuthService/internal/SSO"
	"AuthService/internal/cConstants"
	"AuthService/internal/model"
	"AuthService/pkg/secure"
	"AuthService/pkg/utils"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	shield *secure.Shield
	db     *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB, shield *secure.Shield) SSO.Repository {
	return &postgresRepository{db: db, shield: shield}
}

func (p postgresRepository) AddClient(params *model.SignUpParams) (*string, *model.CodeModel) {
	var uuid string
	err := p.db.QueryRow(`INSERT INTO client DEFAULT VALUES
    RETURNING client_uuid`).Scan(&uuid)
	tx, err := p.db.Begin()
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoGetClient)
		return nil, cErr
	}
	stmts := []string{
		"INSERT INTO client_contact (client_uuid, email) VALUES ($1, $2)",
		"INSERT INTO ver_level (client_uuid) VALUES ($1)",
		"INSERT INTO credential (client_uuid, password) VALUES ($1, $2)",
		"INSERT INTO history_passwords (client_uuid, password) VALUES ($1, $2)",
	}

	_, err = tx.Exec(stmts[0], uuid, params.Email)
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoGetClient)
		return nil, cErr
	}
	_, err = tx.Exec(stmts[1], uuid)
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoGetClient)
		return nil, cErr
	}
	_, err = tx.Exec(stmts[2], uuid, params.Password)
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoGetClient)
		return nil, cErr
	}
	_, err = tx.Exec(stmts[3], uuid, params.Password)
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoGetClient)
		return nil, cErr
	}
	err = tx.Commit()
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoGetClient)
		return nil, cErr
	}
	cErr := p.errorCheckInsert(nil, err, cConstants.RepoGetClient)
	return &uuid, cErr
}

func (p postgresRepository) GetClient(clientUuid string) (*model.Client, *model.CodeModel) {
	var clientRepo []model.Client
	var contacts []model.ClientContacts
	var credential []model.Credential
	var verLevel []model.AuthLevel

	err := p.db.Select(&clientRepo, `SELECT *
									FROM client
									WHERE client_uuid=$1`, clientUuid)
	err = p.db.Select(&contacts, `SELECT *
									FROM client_contact
									WHERE client_uuid=$1`, clientUuid)
	err = p.db.Select(&credential, `SELECT *
									FROM credential
									WHERE client_uuid=$1`, clientUuid)
	err = p.db.Select(&verLevel, `SELECT *
									FROM ver_level
									WHERE client_uuid=$1`, clientUuid)
	cErr := p.errorCheckSelect(clientRepo, err, cConstants.RepoGetClient)
	if cErr.IsError {
		return nil, cErr
	}
	if clientRepo == nil {
		return nil, cErr
	}
	if len(clientRepo) == 0 {
		return nil, &model.CodeModel{
			InternalMessage: "no such client",
		}
	}
	if len(contacts) == 0 {
		return nil, &model.CodeModel{
			InternalMessage: "no such contact",
		}
	}
	if len(credential) == 0 {
		return nil, &model.CodeModel{
			InternalMessage: "no such credential",
		}
	}
	if len(verLevel) == 0 {
		return nil, &model.CodeModel{
			InternalMessage: "no such verLevel",
		}
	}
	client := &clientRepo[0]
	client.Contacts = &contacts[0]
	client.Credential = &credential[0]
	client.AuthLevel = &verLevel[0]
	return client, cErr
}

func (p postgresRepository) GetClientByTG(params *model.SignInTGParams) (*model.Client, *model.CodeModel) {
	var data []model.Client
	err := p.db.Select(&data, `SELECT clt.*, cont.*, vl.*, cred.*
									FROM client clt
									INNER JOIN client_contact cont
									on cont.client_uuid = clt.client_uuid
									INNER JOIN ver_level vl
									on cont.client_uuid = clt.client_uuid
									INNER JOIN credential cred
									on cont.client_uuid = clt.client_uuid
									WHERE cont.tg=$1`, params.TgUserId)

	cErr := p.errorCheckSelect(data, err, cConstants.RepoGetClient)
	if cErr.IsError {
		return nil, cErr
	}
	if data == nil {
		return nil, cErr
	}
	return &data[0], cErr
}

func (p postgresRepository) GetClientAuthLevel(params *string) (*model.AuthLevel, *model.CodeModel) {
	var data []model.AuthLevel
	err := p.db.Select(&data, `SELECT * FROM ver_level WHERE client_uuid=$1`, params)
	cErr := p.errorCheckSelect(data, err, cConstants.RepoGetClientAuthLevel)
	if data == nil {
		return nil, cErr
	}
	return &data[0], cErr
}

//--------------------------------------------------------------------------------------

func (p postgresRepository) GetClientUuidByEmail(email string) (*string, *model.CodeModel) {
	q := `SELECT client_uuid FROM client_contact WHERE LOWER(email)=$1`
	var data []model.ClientUuid
	err := p.db.Select(&data, q, email)
	cErr := p.errorCheckSelect(data, err, cConstants.RepoGetClientUuidByEmail)
	if data == nil {
		return nil, cErr
	}
	return &data[0].ClientUuid, cErr
}

//--------------------------------------------------------------------------------------

func (p postgresRepository) GetNicknameHistory(params *string) (*[]model.NicknameHistory, *model.CodeModel) {
	q := `SELECT * FROM history_nickname WHERE client_uuid=$2`
	var data []model.NicknameHistory
	stmt, err := p.db.Preparex(q)
	err = stmt.Select(data, params)
	cErr := p.errorCheckSelect(data, err, cConstants.RepoGetNicknameHistory)
	return &data, cErr
}

func (p postgresRepository) WriteNicknameHistory(params *model.NicknameHistory) *model.CodeModel {
	q := `INSERT INTO history_nickname (client_uuid, old_nickname) VALUES ($1, $2)`
	stmt, err := p.db.Preparex(q)
	res, err := stmt.Exec(params, params.OldNickname)
	cErr := p.errorCheckInsert(res, err, cConstants.RepoChangeLevelStatus)
	return cErr
}

//--------------------------------------------------------------------------------------

func (p postgresRepository) GetCredential(params *string) (*model.Credential, *model.CodeModel) {
	q := `SELECT * FROM credential WHERE client_uuid=$1`
	var data []model.Credential
	stmt, err := p.db.Preparex(q)
	err = stmt.Select(&data, params)
	cErr := p.errorCheckSelect(data, err, cConstants.RepoGetCredential)
	if data == nil {
		return nil, cErr
	}
	return &data[0], cErr
}

//--------------------------------------------------------------------------------------

func (p postgresRepository) ChangePassword(params *model.ChangePasswordParams) *model.CodeModel {
	params.NewPassword, _ = p.encode(params.NewPassword)
	res, err := p.db.Exec(`UPDATE credential
								  SET password=$1 WHERE client_uuid=$2`, params.NewPassword, params.ClientUuid)

	cErr := p.errorCheckUpdate(res, err, cConstants.RepoGetClient)
	return cErr
}

func (p postgresRepository) GetPassHistory(clientUuid *string) (*[]model.HistoryPasswords, *model.CodeModel) {
	var data []model.HistoryPasswords
	err := p.db.Select(&data, `SELECT * FROM history_passwords
									WHERE client_uuid=$1
									`, &clientUuid)

	cErr := p.errorCheckSelect(data, err, cConstants.RepoGetClient)
	if cErr.IsError {
		return nil, cErr
	}
	return &data, cErr
}

//--------------------------------------------------------------------------------------

func (p postgresRepository) GetCodeRecovery(params *model.GetCodeRecoveryParams) (*[]model.AuthCode, *model.CodeModel) {
	var data []model.AuthCode

	q := `SELECT code_need, destination, client_uuid FROM auth_db.public.codes_confirms
	WHERE code_need=$1`
	err := p.db.Select(&data, q, params.Hash)

	cErr := p.errorCheckSelect(data, err, cConstants.RepoGetAuthCodeByType)
	if cErr.IsError {
		return nil, cErr
	}
	if data == nil {
		return nil, cErr
	}
	return &data, cErr
}

func (p postgresRepository) GetAuthCodeByType(params *model.GetAuthCodeByTypeParams) (*[]model.AuthCode, *model.CodeModel) {
	var data []model.AuthCode

	q := `SELECT code_need, client_uuid, destination 
		  FROM auth_db.public.codes_confirms
		  WHERE client_uuid=$1 AND type_code=$2`
	q += fmt.Sprintf(" AND create_time > now() + interval '3 hours' - interval '%s'", cConstants.CodesConfirmExpiredTime[params.Type])
	q += " ORDER BY create_time DESC"
	err := p.db.Select(&data, q, params.ClientUuid, params.Type)

	cErr := p.errorCheckSelect(data, err, cConstants.RepoGetAuthCodeByType)
	if cErr.IsError {
		return nil, cErr
	}
	return &data, cErr
}

func (p postgresRepository) WriteAuthCode(params *model.WriteAuthCodeParams) *model.CodeModel {
	res, err := p.db.NamedExec(`INSERT INTO auth_db.public.codes_confirms (type_code, code_need, client_uuid, destination)
												VALUES (:type_code, :code_need, :client_uuid, :destination)`, params.Code)
	cErr := p.errorCheckInsert(res, err, cConstants.RepoWriteAuthCode)
	return cErr
}

//--------------------------------------------------------------------------------------

func (p postgresRepository) ChangeLevelStatus(params *model.AuthLevelUpdateStatusParams) *model.CodeModel {
	q := fmt.Sprintf("UPDATE ver_level SET %s=$1 WHERE client_uuid=$2", params.LevelName)
	fmt.Println(q, params.ClientUuid)
	stmt, err := p.db.Prepare(q)
	res, err := stmt.Exec(params.IsValid, params.ClientUuid)
	cErr := p.errorCheckUpdate(res, err, cConstants.RepoChangeLevelStatus)
	return cErr
}

func (p postgresRepository) AddTotpSecret(params *model.AddTotpSecretParams) *model.CodeModel {
	q := `UPDATE credential SET totp_secret=$1 WHERE client_uuid=$2`
	stmt, err := p.db.Prepare(q)
	res, err := stmt.Exec(params.SecretTotp, params.ClientUuid)
	cErr := p.errorCheckInsert(res, err, cConstants.RepoChangeLevelStatus)
	return cErr
}

//--------------------------------------------------------------------------------------

func (p postgresRepository) ChangePhone(params *model.ChangePhoneParams) *model.CodeModel {
	q := `UPDATE client_contact SET phone=$1 WHERE client_uuid=$2`
	stmt, err := p.db.Prepare(q)
	res, err := stmt.Exec(params.NewPhone, params.ClientUuid)
	cErr := p.errorCheckInsert(res, err, cConstants.RepoChangeLevelStatus)
	return cErr
}

func (p postgresRepository) ChangeTg(params *model.ChangeTgParams) *model.CodeModel {
	stmts := []string{
		"UPDATE client_contact SET tg=$1 WHERE client_uuid=$2",
		"UPDATE ver_level SET tg=true WHERE client_uuid=$1;",
	}
	tx, err := p.db.Begin()
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoChangeTg)
		return cErr
	}
	_, err = tx.Exec(stmts[0], params.NewTg, params.ClientUuid)
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoChangeTg)
		return cErr
	}
	_, err = tx.Exec(stmts[1], params.ClientUuid)
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoChangeTg)
		return cErr
	}
	err = tx.Commit()
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoChangeTg)
		return cErr
	}
	cErr := p.errorCheckInsert(nil, err, cConstants.RepoChangeTg)
	return cErr
}

func (p postgresRepository) ChangeTgID(params *model.ChangeTgParams) *model.CodeModel {
	q := `  BEGIN;
					UPDATE credential SET tg_id=$1 
					WHERE client_uuid=$2;
					
					UPDATE ver_level SET tg=true
					WHERE client_uuid=$2;
			COMMIT;`
	stmt, err := p.db.Prepare(q)
	res, err := stmt.Exec(params.NewTg, params)
	cErr := p.errorCheckInsert(res, err, cConstants.RepoChangeLevelStatus)
	return cErr
}

func (p postgresRepository) ChangeNickname(params *model.ChangeNicknameParams) *model.CodeModel {
	tx, err := p.db.Begin()
	stmts := []string{
		"UPDATE client  SET nickname=$1  WHERE client_uuid=$2;",
		"INSERT INTO history_nickname (client_uuid, old_nickname)   VALUES ($1, $2);",
	}
	_, err = tx.Exec(stmts[0], params.NewNickname, params.ClientUuid)
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoChangeNickname)
		return cErr
	}
	_, err = tx.Exec(stmts[1], params.ClientUuid, params.OldNickname)
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoChangeNickname)
		return cErr
	}
	err = tx.Commit()
	if err != nil {
		cErr := p.errorCheckInsert(nil, err, cConstants.RepoChangeNickname)
		return cErr
	}

	cErr := p.errorCheckInsert(nil, err, cConstants.RepoChangeNickname)
	return cErr
}

//--------------------------------------------------------------------------------------

func (p postgresRepository) GetActiveSessions(clientUuid *string) (*[]model.Session, *model.CodeModel) {
	q := `SELECT * FROM session_history WHERE client_uuid=$1 AND is_logout=false`
	var data []model.Session
	stmt, err := p.db.Preparex(q)
	err = stmt.Select(&data, clientUuid)
	cErr := p.errorCheckSelect(data, err, cConstants.RepoChangeLevelStatus)
	return &data, cErr
}

func (p postgresRepository) AddSession(params *model.Session) *model.CodeModel {
	q := `INSERT INTO session_history
	(client_uuid, ip, ua) VALUES ($1, $2, $3)`
	stmt, err := p.db.Preparex(q)
	res, err := stmt.Exec(params.ClientUuid, params.IP, params.UA)
	cErr := p.errorCheckInsert(res, err, cConstants.RepoChangeLevelStatus)
	return cErr
}

func (p postgresRepository) LogoutSession(params *model.LogoutParams) *model.CodeModel {
	q := `UPDATE session_history SET logout_time=$1, is_logout=true  WHERE client_uuid=$2`
	stmt, err := p.db.Prepare(q)
	res, err := stmt.Exec(utils.GetEuropeTime(), params.ClientUuid)
	cErr := p.errorCheckInsert(res, err, cConstants.RepoChangeLevelStatus)
	return cErr
}

//--------------------------------------------------------------------------------------

func (p postgresRepository) errorCheckSelect(data interface{}, err error, interCode model.OmitZeroInt) *model.CodeModel {
	var cm = &model.CodeModel{}
	if err != nil {
		cm = model.GetError(interCode, fiber.StatusInternalServerError, err.Error(), err.Error())
		return cm
	}

	if utils.GetInterfaceLength(data) == 0 {
		cm = model.GetError(interCode, fiber.StatusBadRequest, cConstants.NoDataRepo, cConstants.NoDataRepo)
	}
	return cm
}

func (p postgresRepository) errorCheckInsert(res sql.Result, err error, interCode model.OmitZeroInt) *model.CodeModel {
	var cm = &model.CodeModel{}
	if err != nil {
		return model.GetError(interCode, cConstants.StatusInternalServerError, err.Error(), err.Error())
	}
	if res == nil {
		return cm
	}
	insertedRow, err := res.RowsAffected()
	if err != nil {
		return model.GetError(interCode, cConstants.StatusInternalServerError, cConstants.RepoInsertedRow, cConstants.RepoInsertedRow)
	}
	if insertedRow == 0 {
		return model.GetError(interCode, cConstants.StatusInternalServerError, cConstants.RepoZeroInsertedRow, cConstants.RepoZeroInsertedRow)
	}
	return cm
}

func (p postgresRepository) errorCheckUpdate(res sql.Result, err error, interCode model.OmitZeroInt) *model.CodeModel {
	var cm = &model.CodeModel{}
	if err != nil {
		return model.GetError(interCode, cConstants.StatusInternalServerError, err.Error(), err.Error())
	}
	if res == nil {
		return cm
	}
	return cm
}

func (p postgresRepository) encode(password string) (string, *model.CodeModel) {
	cryptPass, err := p.shield.EncryptMessage(password)
	if err != nil {
		return "", model.GetError(cConstants.RepoErrToDecodePassword, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return cryptPass, &model.CodeModel{}
}

func (p postgresRepository) decode(password string) string {
	cryptPass := p.shield.DecryptMessage(password)

	return cryptPass
}
