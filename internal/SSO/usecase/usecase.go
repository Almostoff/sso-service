package usecase

import (
	"AuthService/config"
	"AuthService/internal/SSO"
	"AuthService/internal/cConstants"
	"AuthService/internal/emailService"
	"AuthService/internal/model"
	"AuthService/internal/redisUsers"
	"AuthService/pkg/logger"
	"AuthService/pkg/secure"
	"AuthService/pkg/utils"
	"bytes"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp/totp"
	"image/png"
	"io/ioutil"
	"log"
	"strings"
)

type UsersUsecase struct {
	logger *logger.ApiLogger
	repo   SSO.Repository
	shield *secure.Shield
	redis  redisUsers.UseCase
	email  emailService.Email
	kyc    config.Kyc
	server config.ServerConfig
}

func NewUsersUsecase(logger *logger.ApiLogger, repo SSO.Repository, shield *secure.Shield, redis redisUsers.UseCase,
	email emailService.Email, kyc config.Kyc, server config.ServerConfig) SSO.UseCase {
	return &UsersUsecase{
		logger: logger,
		repo:   repo,
		shield: shield,
		redis:  redis,
		email:  email,
		kyc:    kyc,
		server: server,
	}
}

func (u *UsersUsecase) SignInTg(params *model.SignInTGParams) *model.ResponseSignIn {
	if params.UA != cConstants.TgUA {
		cErr := model.GetError(cConstants.UCUaNotTG, cConstants.StatusInternalServerError, cConstants.UaNotTG, cConstants.WrongCred)
		return &model.ResponseSignIn{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	client, cErr := u.repo.GetClientByTG(params)
	if cErr.IsError {
		return &model.ResponseSignIn{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}

	session := model.Session{ClientUuid: client.ClientUuid, UA: params.UA, IP: cConstants.TGIP}
	data, cErr := u.AddRedisUser(&session, client.Contacts.Email)
	if cErr.IsError {
		return &model.ResponseSignIn{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	cErr = u.repo.AddSession(&session)
	if cErr.IsError {
		return &model.ResponseSignIn{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	return &model.ResponseSignIn{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
		Data:    data,
	}

}

func (u *UsersUsecase) SignUp(params *model.SignUpParams) *model.ResponseSignUp {
	var cErr *model.CodeModel
	params.Password = u.encode(params.Password)
	email := strings.ToLower(params.Email)
	uuid, err := u.GetClientUuidByEmail(email)
	if err.IsError {
		return &model.ResponseSignUp{
			Error:   err,
			Success: model.GetSuccessByCM(err),
		}
	}
	if uuid == nil {
		uuid, cErr = u.repo.AddClient(params)
		if cErr.IsError {
			return &model.ResponseSignUp{
				Error:   cErr,
				Success: model.GetSuccessByCM(cErr),
			}
		}
	}

	session := &model.Session{ClientUuid: *uuid, UA: params.UA, IP: params.IP}
	//data, cErr := u.AddRedisUser(session, params.Email)
	//if cErr.IsError {
	//	return &model.ResponseSignUp{
	//		Error:   cErr,
	//		Success: model.GetSuccessByCM(cErr),
	//	}
	//}
	data, cErr := u.AddRedisUser(session, email)
	if cErr.IsError {
		return &model.ResponseSignUp{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	cErr = u.repo.AddSession(session)
	if cErr.IsError {
		return &model.ResponseSignUp{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	return &model.ResponseSignUp{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
		Data:    (*model.SignUpData)(data),
	}
}

func (u *UsersUsecase) SignIn(params *model.SignInParams) *model.ResponseSignIn {
	email := strings.ToLower(params.Email)
	uuid, cErr := u.repo.GetClientUuidByEmail(email)
	if cErr.IsError {
		return &model.ResponseSignIn{
			Error:   cErr,
			Success: model.GetSuccess(false, ""),
		}
	}
	user := u.GetClient(*uuid)
	if user.Error.InternalCode != 0 {
		return &model.ResponseSignIn{
			Error:   user.Error,
			Success: model.GetSuccess(false, ""),
		}
	}
	pass := u.decode(user.Data.Credential.Password)
	if pass != params.Password {
		cErr := model.GetError(cConstants.RepoErrToDecodePassword, cConstants.StatusUnauthorized, "", "")
		return &model.ResponseSignIn{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	session := model.Session{ClientUuid: params.ClientUuid, UA: params.UA, IP: params.IP}
	data, cErr := u.AddRedisUser(&session, email)
	if cErr.IsError {
		return &model.ResponseSignIn{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	cErr = u.repo.AddSession(&session)
	if cErr.IsError {
		return &model.ResponseSignIn{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	return &model.ResponseSignIn{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
		Data:    data,
	}

}

func (u *UsersUsecase) GetClient(params string) *model.ResponseClient {
	client, cErr := u.repo.GetClient(params)
	return &model.ResponseClient{
		Success: model.GetSuccess(true, ""),
		Error:   cErr,
		Data:    client,
	}
}

func (u *UsersUsecase) GetClientPublicInfo(params string) *model.ResponseClient {
	client, cErr := u.repo.GetClient(params)
	client.Credential = nil
	client.Contacts = nil
	return &model.ResponseClient{
		Success: model.GetSuccess(true, ""),
		Error:   cErr,
		Data:    client,
	}
}

func (u *UsersUsecase) GetClientPrivateInfo(clientUuid string) *model.ResponseClient {
	client, cErr := u.repo.GetClient(clientUuid)
	client.Credential = nil
	return &model.ResponseClient{
		Success: model.GetSuccess(true, ""),
		Error:   cErr,
		Data:    client,
	}
}

// ---------------------------------------------------------------------

func (u *UsersUsecase) ConfirmMail(params *model.ConfirmMailParams) *model.Response {
	valid := u.IsCodeValid(&model.VerCodeParams{
		CodeInput: params.Hash,
		Type:      cConstants.Email,
	})
	if !valid.Success.Success || valid.Error.IsError {
		return valid
	}
	cErr := u.repo.ChangeLevelStatus(&model.AuthLevelUpdateStatusParams{
		ClientUuid: valid.Success.Message,
		LevelName:  cConstants.Em,
		IsValid:    true,
	})
	return &model.Response{
		Success: model.GetSuccessByCM(cErr),
		Error:   cErr,
	}
}

func (u *UsersUsecase) RequestToConfirmMail(params *model.RequestToConfirmMailParams) *model.Response {
	client := u.GetClient(params.ClientUuid)
	if client.Error.IsError {
		return &model.Response{
			Error:   client.Error,
			Success: model.GetSuccessByCM(client.Error),
		}
	}
	code := utils.GenerateHash()
	res := u.AddCode(&model.WriteAuthCodeParams{
		ClientUuid: params.ClientUuid,
		Code: &model.AuthCode{
			ClientUuid:  params.ClientUuid,
			CodeNeed:    code,
			Type:        cConstants.Email,
			Destination: client.Data.Contacts.Email,
		}})
	go u.email.SendMailConfirm(&emailService.SendMailConfirmParams{
		TypeEmail:   cConstants.ConfirmEmail,
		Link:        code,
		LanguageIso: client.Data.Language,
		Email:       client.Data.Contacts.Email,
	})
	return res
}

// ---------------------------------------------------------------------

func (u *UsersUsecase) Logout(params *model.LogoutParams) *model.Response {
	log.Println("Check params client ID: ", params.ClientUuid)
	cErr := u.DeleteRedisUser(&model.DeleteRedisUserParams{
		ClientUuid: params.ClientUuid,
	})
	if cErr.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccess(false, cConstants.WrongCred),
		}
	}
	cErr = u.repo.LogoutSession(params)
	if cErr.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccess(true, ""),
	}
}

func (u *UsersUsecase) GetActiveSession(params *model.GetActiveSessionParams) *model.ResponseGetActiveSession {
	data, cErr := u.repo.GetActiveSessions(&params.ClientUuid)
	return &model.ResponseGetActiveSession{
		Error:   cErr,
		Data:    data,
		Success: model.GetSuccessByCM(cErr),
	}
}

// ---------------------------------------------------------------------

func (u *UsersUsecase) GetAuthLevel(params *string) *model.ResponseAuthLevel {
	level, cErr := u.repo.GetClientAuthLevel(params)
	return &model.ResponseAuthLevel{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
		Data:    level,
	}
}

// ---------------------------------------------------------------------

func (u *UsersUsecase) ChangePassword(params *model.ChangePasswordParams) *model.Response {
	history, cErr := u.repo.GetPassHistory(&params.ClientUuid)
	if cErr.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccess(false, ""),
		}
	}
	in := u.checkIfPass(*history, params.NewPassword)
	if in {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccess(false, cConstants.NeedNewPass),
		}
	}
	cErr = u.repo.ChangePassword(&model.ChangePasswordParams{
		ClientUuid:  params.ClientUuid,
		NewPassword: params.NewPassword,
	})
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccess(true, ""),
	}
}

func (u *UsersUsecase) ChangePasswordWithOldCheck(params *model.ChangePasswordParams) *model.Response {
	cred, cErr := u.repo.GetCredential(&params.ClientUuid)
	if cErr.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccess(false, ""),
		}
	}
	pass := u.decode(cred.Password)
	if pass != params.OldPassword {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccess(false, cConstants.WrongCred),
		}
	}
	res := u.ChangePassword(params)
	if res.Error.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccess(false, ""),
		}
	}
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccess(true, ""),
	}
}

func (u *UsersUsecase) ChangePhone(params *model.ChangePhoneParams) *model.Response {
	cErr := u.repo.ChangePhone(params)
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
	}
}

func (u *UsersUsecase) ChangeTg(params *model.ChangeTgParams) *model.Response {
	cErr := u.repo.ChangeTg(params)
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
	}
}

func (u *UsersUsecase) ChangeNickname(params *model.ChangeNicknameParams) *model.Response {
	client, cErr := u.repo.GetClient(params.ClientUuid)
	if cErr.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	params.OldNickname = client.Nickname
	cErr = u.repo.ChangeNickname(params)
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
	}
}

// ---------------------------------------------------------------------

func (u *UsersUsecase) AddTotp(params *model.AddTotpParams) *model.ResponseAddTotp {
	client, cErr := u.repo.GetClient(params.ClientUuid)
	if cErr.InternalCode != 0 {
		return &model.ResponseAddTotp{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	var accountName string
	if client.Nickname == "" {
		accountName = client.Contacts.Email
	} else {
		accountName = client.Nickname
	}
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "p2p.exnode.ru",
		AccountName: accountName,
		SecretSize:  50,
	})
	if err != nil {
		return &model.ResponseAddTotp{
			Success: model.GetSuccessByCM(cErr),
			Error:   model.GetError(cConstants.GenerateTotp, cConstants.StatusInternalServerError, err.Error(), err.Error()),
		}
	}
	var buf bytes.Buffer
	img, err := secret.Image(200, 200)
	err = png.Encode(&buf, img)
	if err != nil {
		return &model.ResponseAddTotp{
			Success: model.GetSuccessByCM(cErr),
			Error:   model.GetError(cConstants.GenerateTotp, cConstants.StatusInternalServerError, err.Error(), err.Error()),
		}
	}
	err = ioutil.WriteFile("qr-code.png", buf.Bytes(), 0644)
	if err != nil {
		return &model.ResponseAddTotp{
			Success: model.GetSuccessByCM(cErr),
			Error:   model.GetError(cConstants.GenerateTotp, cConstants.StatusInternalServerError, err.Error(), err.Error()),
		}
	}

	cErr = u.repo.AddTotpSecret(&model.AddTotpSecretParams{
		SecretTotp: secret.Secret(),
		ClientUuid: params.ClientUuid})
	if cErr.InternalCode != 0 {
		return &model.ResponseAddTotp{
			Success: model.GetSuccessByCM(cErr),
			Error:   cErr,
		}
	}
	return &model.ResponseAddTotp{
		Success: model.GetSuccessByCM(cErr),
		Data: &model.TotpModel{
			File:        buf.Bytes(),
			AccountName: accountName,
			Secret:      secret.Secret(),
			Link:        secret.String(),
		},
		Error: cErr,
	}

}

func (u *UsersUsecase) VerifyTotpInit(params *model.VerifyTotpParams) *model.Response {
	ok := u.VerifyTotp(params)
	if !ok.Success.Success {
		return ok
	}
	cErr := u.repo.ChangeLevelStatus(&model.AuthLevelUpdateStatusParams{ClientUuid: params.ClientUuid, LevelName: cConstants.Totp, IsValid: true})
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
	}
}

func (u *UsersUsecase) VerifyTotp(params *model.VerifyTotpParams) *model.Response {
	client, cErr := u.repo.GetClient(params.ClientUuid)
	if cErr.InternalCode != 0 {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccess(false, ""),
		}
	}
	ok := totp.Validate(params.Token, client.Credential.TotpSecret)
	if !ok {
		return &model.Response{
			Error:   model.GetError(cConstants.WrongTotp, cConstants.StatusBadRequest, cConstants.WrongCred, cConstants.WrongCred),
			Success: model.GetSuccess(false, cConstants.WrongCred),
		}
	}
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
	}

}

// ---------------------------------------------------------------------

func (u *UsersUsecase) AddCode(params *model.WriteAuthCodeParams) *model.Response {
	params.Code.Destination = strings.ToLower(params.Code.Destination)
	cErr := u.repo.WriteAuthCode(params)
	return &model.Response{
		Success: model.GetSuccessByCM(cErr),
		Error:   cErr,
	}
}

func (u *UsersUsecase) IsCodeValid(params *model.VerCodeParams) *model.Response {
	if params.Type == cConstants.RecoveryAccess || params.Type == cConstants.Email {
		codes, cErr := u.repo.GetCodeRecovery(&model.GetCodeRecoveryParams{
			Hash: params.CodeInput,
		})
		if cErr.IsError {
			return &model.Response{
				Error:   cErr,
				Success: model.GetSuccessByCM(cErr),
			}
		}
		var ok bool
		var clientUuid string
		for _, v := range *codes {
			if v.CodeNeed == params.CodeInput {
				ok = true
				clientUuid = v.ClientUuid
			}
		}
		return &model.Response{
			Success: model.GetSuccess(ok, clientUuid),
			Error:   cErr,
		}

	} else {
		codes, cErr := u.repo.GetAuthCodeByType(&model.GetAuthCodeByTypeParams{
			Type:       params.Type,
			ClientUuid: params.ClientUuid,
		})
		if cErr.IsError {
			return &model.Response{
				Error:   cErr,
				Success: model.GetSuccessByCM(cErr),
			}
		}
		var ok bool
		var clientUuid string
		for _, v := range *codes {
			if v.CodeNeed == params.CodeInput {
				ok = true
				clientUuid = v.ClientUuid
			}
		}
		return &model.Response{
			Success: model.GetSuccess(ok, clientUuid),
			Error:   cErr,
		}
	}

}

// ---------------------------------------------------------------------

func (u *UsersUsecase) RecoveryInit(params *model.RecoveryInitParams) *model.Response {
	clientUuid, cErr := u.GetClientUuidByEmail(params.Email)
	if cErr.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	client := u.GetClient(*clientUuid)
	if client.Error.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	code := utils.GenerateHash()
	email := strings.ToLower(client.Data.Contacts.Email)
	cErr = u.repo.WriteAuthCode(&model.WriteAuthCodeParams{
		ClientUuid: *clientUuid,
		Code: &model.AuthCode{
			ClientUuid:  *clientUuid,
			CodeNeed:    code,
			Destination: email,
			Type:        cConstants.RecoveryByEmail,
		}})
	go u.email.SendRecoveryLink(&emailService.SendRecoveryLinkParams{
		Link:        code,
		Email:       client.Data.Contacts.Email,
		LanguageIso: client.Data.Language,
		TypeEmail:   cConstants.RecoveryAccess,
	})
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
	}
}

func (u *UsersUsecase) RecoveryConfirm(params *model.RecoveryConfirmParams) *model.Response {
	codesList, cErr := u.repo.GetCodeRecovery(&model.GetCodeRecoveryParams{
		Hash: params.Hash,
	})
	if cErr.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	codes := *codesList
	code := codes[0]
	if code.CodeNeed != params.Hash {
		return &model.Response{
			Error:   model.GetError(cConstants.WrongHashRecovery, cConstants.StatusBadRequest, cConstants.WrongCred, cConstants.WrongCred),
			Success: model.GetSuccessByCM(cErr),
		}
	}
	res := u.ChangePassword(&model.ChangePasswordParams{
		ClientUuid:  code.ClientUuid,
		NewPassword: params.Password,
	})
	if res.Error.IsError || !res.Success.Success {
		return &model.Response{
			Error:   model.GetError(cConstants.RecoveryConfirmPass, cConstants.StatusBadRequest, cConstants.NeedNewPass, cConstants.NeedNewPass),
			Success: model.GetSuccessByCM(cErr),
		}
	}
	return &model.Response{
		Error:   res.Error,
		Success: model.GetSuccessByCM(res.Error),
	}
}

// ---------------------------------------------------------------------

func (u *UsersUsecase) AddRedisUser(params *model.Session, email string) (*model.SignInData, *model.CodeModel) {
	accessToken, err := u.generateAccessToken(&model.TokenGenerate{ClientUuid: params.ClientUuid, Email: email})
	if err != nil {
		return nil, model.GetError(cConstants.GenerateAccess, cConstants.StatusInternalServerError, err.Error(), err.Error())
	}
	refreshToken, err := u.generateRefreshToken(&model.TokenGenerate{ClientUuid: params.ClientUuid, Email: email})
	if err != nil {
		return nil, model.GetError(cConstants.GenerateRefresh, cConstants.StatusInternalServerError, err.Error(), err.Error())
	}
	return &model.SignInData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, &model.CodeModel{}
}

func (u *UsersUsecase) DeleteRedisUser(params *model.DeleteRedisUserParams) *model.CodeModel {
	fingerprint := fmt.Sprintf("%s%s", params.UA, params.ClientUuid)
	deletedSessions, err := u.redis.RemoveUser(&redisUsers.RemoveUser{
		Fingerprint: fingerprint,
		Uuid:        params.ClientUuid,
	})
	if err != nil {
		return &model.CodeModel{
			InternalMessage: err.Error(),
		}
	}
	if deletedSessions == 0 {
		return &model.CodeModel{
			InternalMessage: "access_token has expired",
		}
	}
	return nil
}

func (u *UsersUsecase) IsValidRedisUser(params *model.LogoutParams) *model.CodeModel {
	return &model.CodeModel{}
}

// ---------------------------------------------------------------------

func (u *UsersUsecase) GetClientUuidByEmail(email string) (*string, *model.CodeModel) {
	uuid, cErr := u.repo.GetClientUuidByEmail(email)
	return uuid, cErr
}

// ---------------------------------------------------------------------

func (u *UsersUsecase) ValidateToken(params *model.ValidateAccessTokenParams) *model.Response {
	cErr := u.IsValidRedisUser(&model.LogoutParams{})
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
	}
}

func (u *UsersUsecase) ValidateAccessToken(params *model.ValidateAccessTokenParams) *model.ResponseValidateToken {
	ok, err := u.redis.Validate(params.Access, fmt.Sprintf("%s%s", params.UA, params.ClientUuid), params.ClientUuid)
	if err != nil {
		return &model.ResponseValidateToken{
			Error: &model.CodeModel{
				InternalCode:    cConstants.StatusInternalServerError,
				InternalMessage: err.Error(),
			},
			Data: &model.ResponseSuccessModel{Success: false},
		}
	}
	return &model.ResponseValidateToken{
		Error: &model.CodeModel{},
		Data:  &model.ResponseSuccessModel{Success: ok},
	}
}

func (u *UsersUsecase) ValidateRefreshToken(params *model.ValidateRefreshTokenParams) *model.ResponseValidateToken {
	ok, err := u.redis.ValidateRefresh(params.Refresh, fmt.Sprintf("%s%s", params.UA, params.ClientUuid), params.ClientUuid)
	if err != nil {
		return &model.ResponseValidateToken{
			Error: &model.CodeModel{
				InternalCode:    cConstants.StatusInternalServerError,
				InternalMessage: err.Error(),
			},
			Data: &model.ResponseSuccessModel{Success: false},
		}
	}
	return &model.ResponseValidateToken{
		Error: &model.CodeModel{},
		Data:  &model.ResponseSuccessModel{Success: ok},
	}
}

func (u *UsersUsecase) RefreshAccessToken(params *model.RefreshAccessTokenParams) *model.ResponseRefreshAccessToken {
	client, cErr := u.repo.GetClient(params.ClientUuid)
	if cErr.InternalCode != 0 {
		return &model.ResponseRefreshAccessToken{
			Error: cErr,
			Data:  &model.Access{},
		}
	}
	accessToken, err := u.generateAccessToken(&model.TokenGenerate{
		ClientUuid: params.ClientUuid,
		Email:      client.Contacts.Email,
	})
	if err != nil {
		return &model.ResponseRefreshAccessToken{
			Data: &model.Access{},
			Error: &model.CodeModel{
				InternalCode:    cConstants.StatusInternalServerError,
				InternalMessage: "access token error",
			},
		}
	}
	newAccess, err := u.redis.UpdateAccess(params.Refresh, accessToken,
		fmt.Sprintf("%s%s", params.UA, params.ClientUuid), client.ClientUuid)
	if err != nil {
		return &model.ResponseRefreshAccessToken{
			Error: &model.CodeModel{
				InternalCode:    cConstants.StatusInternalServerError,
				InternalMessage: err.Error(),
			},
			Data: &model.Access{},
		}
	}
	if newAccess == "" {
		return &model.ResponseRefreshAccessToken{
			Error: &model.CodeModel{
				InternalCode:    cConstants.StatusInternalServerError,
				InternalMessage: "not valid refresh",
			},
			Data: &model.Access{},
		}
	}
	return &model.ResponseRefreshAccessToken{
		Error: &model.CodeModel{},
		Data:  &model.Access{Access: newAccess},
	}
}

// ---------------------------------------------------------------------

func (u *UsersUsecase) KycConfirmInit(params *model.KycConfirmInitParams) *model.Response {
	verLevel, cErr := u.repo.GetClientAuthLevel(&params.ClientUuid)
	if cErr.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	if verLevel.KYC {
		return &model.Response{
			Success: model.GetSuccess(false, cConstants.AlreadyHaveKyc),
		}
	}
	codesList, cErr := u.repo.GetAuthCodeByType(&model.GetAuthCodeByTypeParams{
		Type:       cConstants.Kyc,
		ClientUuid: params.ClientUuid,
	})
	codes := *codesList
	t1 := codes[0].Date
	t2 := utils.GetEuropeTime()
	delta := t2.Sub(t1).Hours()
	if delta < 24 {
		return &model.Response{
			Success: model.GetSuccess(false, "please wait 24 hours"),
		}
	}

	hashID, _ := u.shield.EncryptMessage(params.ClientUuid)
	cErr = u.repo.WriteAuthCode(&model.WriteAuthCodeParams{
		ClientUuid: params.ClientUuid,
		Code: &model.AuthCode{
			ClientUuid: params.ClientUuid,
			CodeNeed:   hashID,
			Type:       cConstants.KYC,
		}})
	if cErr.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}

	var responseModel model.ResponseConfirmKycModel
	http := resty.New().EnableTrace().SetDebug(true)
	var body model.HashKycConfirmBody
	callbackUrl := fmt.Sprintf("http://%v:%v/sso/my/confirm/verify/kyc", u.server.Host, u.server.Port)
	body.CallBack = callbackUrl
	body.ClientID = hashID
	response, err := http.R().SetBasicAuth(u.kyc.Username, u.kyc.Pass).SetBody(body).SetResult(&responseModel).Post(u.kyc.Url)
	if err != nil {
		return &model.Response{
			Error:   model.GetError(cConstants.KycConfirmInit, cConstants.StatusInternalServerError, err.Error(), err.Error()),
			Success: model.GetSuccessByCM(cErr),
		}
	}
	if response == nil {
		return &model.Response{
			Error:   model.GetError(cConstants.KycConfirmInit, cConstants.StatusInternalServerError, cConstants.ResponseNil, err.Error()),
			Success: model.GetSuccessByCM(cErr),
		}
	}

	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccess(true, ""),
	}
}

func (u *UsersUsecase) KycConfirm(params *model.KycConfirmParams) *model.Response {
	decryptUuid := u.shield.DecryptMessage(params.Hash)
	codes, cErr := u.repo.GetAuthCodeByType(&model.GetAuthCodeByTypeParams{
		Type:       cConstants.KYC,
		ClientUuid: decryptUuid,
	})
	if cErr.IsError {
		return &model.Response{
			Error:   cErr,
			Success: model.GetSuccessByCM(cErr),
		}
	}
	code := *codes
	uuid := code[0].ClientUuid
	if decryptUuid != uuid {
		return &model.Response{
			Error:   model.GetError(cConstants.StatusBadRequest, cConstants.StatusBadRequest, "error decrypt hash", ""),
			Success: model.GetSuccessByCM(cErr),
		}
	}
	cErr = u.repo.ChangeLevelStatus(&model.AuthLevelUpdateStatusParams{
		ClientUuid: uuid,
		LevelName:  cConstants.Kyc,
	})
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
	}
}

// ---------------------------------------------------------------------

func (u *UsersUsecase) decode(password string) string {
	cryptPass := u.shield.DecryptMessage(password)

	return cryptPass
}

func (u *UsersUsecase) encode(password string) string {
	cryptPass, _ := u.shield.EncryptMessage(password)

	return cryptPass
}

func (u *UsersUsecase) generateAccessToken(params *model.TokenGenerate) (string, error) {
	var hmacSampleSecret []byte
	fmt.Println(params.Email)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": params.Email,
		"uuid":  params.ClientUuid,
		"type":  "access",
		"nbf":   utils.GetEuropeTime(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)

	return tokenString, err
}

func (u *UsersUsecase) generateRefreshToken(params *model.TokenGenerate) (string, error) {
	var hmacSampleSecret []byte

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": params.ClientUuid,
		"type": "refresh",
		"nbf":  utils.GetEuropeTime(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)

	return tokenString, err
}

func (u *UsersUsecase) DecodeToken(token string) string {
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	var email string
	email = fmt.Sprint(claims["email"])
	return email
}

func (u *UsersUsecase) getBadResponse(message string, cErr *model.CodeModel) *model.Response {
	return &model.Response{
		Error: cErr,
		Success: &model.Success{
			Success: false,
			Message: message,
		},
	}
}

func (u *UsersUsecase) decodeListPasswords(hp []model.HistoryPasswords, newPassword string) bool {
	var res bool
	for _, v := range hp {
		if u.decode(v.Password) == newPassword {
			res = true
		}
	}
	return res
}

func (u *UsersUsecase) checkIfPass(hp []model.HistoryPasswords, newPass string) bool {
	for _, p := range hp {
		if u.decode(p.Password) == newPass {
			return true
		}
	}
	return false
}

func (u *UsersUsecase) handleCodes(codes *[]model.AuthCode) *model.Response {
	code := *codes
	uuid := code[0].ClientUuid

	cErr := u.repo.ChangeLevelStatus(&model.AuthLevelUpdateStatusParams{ClientUuid: uuid, LevelName: cConstants.Kyc})
	return &model.Response{
		Error:   cErr,
		Success: model.GetSuccessByCM(cErr),
	}
}
