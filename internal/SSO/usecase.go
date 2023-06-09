package SSO

import "AuthService/internal/model"

type UseCase interface {
	SignIn(params *model.SignInParams) *model.ResponseSignIn
	SignInTg(params *model.SignInTGParams) *model.ResponseSignIn
	SignUp(params *model.SignUpParams) *model.ResponseSignUp

	Logout(params *model.LogoutParams) *model.Response

	AddRedisUser(params *model.Session, email string) (*model.SignInData, *model.CodeModel)
	DeleteRedisUser(params *model.DeleteRedisUserParams) *model.CodeModel

	ConfirmMail(params *model.ConfirmMailParams) *model.Response
	RequestToConfirmMail(params *model.RequestToConfirmMailParams) *model.Response
	//ConfirmPhone(params *model.ConfirmPhoneParams) *model.Response
	//RequestToConfirmPhone(params *model.RequestToConfirmPhoneParams) *model.Response

	ChangePassword(params *model.ChangePasswordParams) *model.Response
	ChangePasswordWithOldCheck(params *model.ChangePasswordParams) *model.Response
	ChangePhone(params *model.ChangePhoneParams) *model.Response
	ChangeTg(params *model.ChangeTgParams) *model.Response
	ChangeNickname(params *model.ChangeNicknameParams) *model.Response

	GetAuthLevel(params *string) *model.ResponseAuthLevel
	GetClientPublicInfo(params string) *model.ResponseClient
	GetClientPrivateInfo(clientUuid string) *model.ResponseClient

	GetActiveSession(params *model.GetActiveSessionParams) *model.ResponseGetActiveSession

	GetClientUuidByEmail(email string) (*string, *model.CodeModel)

	AddTotp(params *model.AddTotpParams) *model.ResponseAddTotp
	VerifyTotpInit(params *model.VerifyTotpParams) *model.Response
	VerifyTotp(params *model.VerifyTotpParams) *model.Response

	RecoveryInit(params *model.RecoveryInitParams) *model.Response
	RecoveryConfirm(params *model.RecoveryConfirmParams) *model.Response

	ValidateToken(params *model.ValidateAccessTokenParams) *model.Response
	ValidateAccessToken(params *model.ValidateAccessTokenParams) *model.ResponseValidateToken
	ValidateRefreshToken(params *model.ValidateRefreshTokenParams) *model.ResponseValidateToken
	RefreshAccessToken(params *model.RefreshAccessTokenParams) *model.ResponseRefreshAccessToken

	KycConfirmInit(params *model.KycConfirmInitParams) *model.Response
	KycConfirm(params *model.KycConfirmParams) *model.Response

	AddCode(params *model.WriteAuthCodeParams) *model.Response
	IsCodeValid(params *model.VerCodeParams) *model.Response

	DecodeToken(token string) string
}
