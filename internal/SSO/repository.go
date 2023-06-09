package SSO

import (
	"AuthService/internal/model"
)

type Repository interface {
	AddClient(params *model.SignUpParams) (*string, *model.CodeModel)
	GetClient(clientUuid string) (*model.Client, *model.CodeModel)
	GetClientAuthLevel(params *string) (*model.AuthLevel, *model.CodeModel)
	GetClientByTG(params *model.SignInTGParams) (*model.Client, *model.CodeModel)
	GetClientUuidByEmail(email string) (*string, *model.CodeModel)

	ChangePassword(params *model.ChangePasswordParams) *model.CodeModel
	ChangePhone(params *model.ChangePhoneParams) *model.CodeModel
	ChangeTg(params *model.ChangeTgParams) *model.CodeModel
	ChangeTgID(params *model.ChangeTgParams) *model.CodeModel
	ChangeNickname(params *model.ChangeNicknameParams) *model.CodeModel
	ChangeLevelStatus(params *model.AuthLevelUpdateStatusParams) *model.CodeModel

	GetCredential(params *string) (*model.Credential, *model.CodeModel)

	GetCodeRecovery(params *model.GetCodeRecoveryParams) (*[]model.AuthCode, *model.CodeModel)
	GetAuthCodeByType(params *model.GetAuthCodeByTypeParams) (*[]model.AuthCode, *model.CodeModel)
	WriteAuthCode(params *model.WriteAuthCodeParams) *model.CodeModel

	GetNicknameHistory(params *string) (*[]model.NicknameHistory, *model.CodeModel)

	GetPassHistory(clientUuid *string) (*[]model.HistoryPasswords, *model.CodeModel)
	AddTotpSecret(params *model.AddTotpSecretParams) *model.CodeModel

	GetActiveSessions(clientUuid *string) (*[]model.Session, *model.CodeModel)
	AddSession(params *model.Session) *model.CodeModel
	LogoutSession(params *model.LogoutParams) *model.CodeModel
}
