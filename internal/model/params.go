package model

type ValidateAccessTokenParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Access     string `json:"Access"`
	UA         string `json:"ua"`
}

type ValidateRefreshTokenParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Refresh    string `json:"Refresh"`
	UA         string `json:"ua"`
}

type AddUserAgentParams struct {
	*Session
}

type LogoutUserAgentParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Access     string `json:"Access"`
	UA         string `json:"ua" db:"ua"`
}

type RefreshAccessTokenParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Refresh    string `json:"Refresh"`
	UA         string `json:"ua"`
}

type GetAuthCodeParams struct {
	Type string `json:"type" db:"type"`
	Hash string
}

type GetAuthCodeByTypeParams struct {
	Type       string `json:"type" db:"type"`
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}

type AuthLevelUpdatePhoneStatusParams struct {
	Hash string `json:"hash" db:"hash"`
}

type LogoutParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Access     string `json:"access"`
	UA         string `json:"ua" db:"ua"`
}

type DeleteRedisUserParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Access     string `json:"access"`
	UA         string `json:"ua" db:"ua"`
}

type AuthLevelUpdateStatusParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	LevelName  string
	IsValid    bool
}

type AuthLevelUpdateKYCStatusParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}

type AuthLevelUpdateTotpStatusParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}

type WriteAuthCodeParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Code       *AuthCode
}

type IsCodeExistParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Type       string `json:"type" db:"type"`
}

type GetCodeRecoveryParams struct {
	Hash string `json:"hash" db:"code_need"`
}

type SignInParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Email      string `json:"email"`
	UA         string `json:"ua" db:"ua"`
	IP         string `json:"ip" db:"ip"`
	Password   string `json:"password" db:"password"`
}

type SignInTGParams struct {
	TgUserName string `json:"tg_user_name"`
	TgUserId   int64  `json:"tg_user_id"`
	UA         string `json:"ua" db:"ua"`
}

type SignInParamsWith2fa struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Password   string `json:"password" db:"password"`
	Auth       ParamsAuth
}

type ParamsAuth struct {
	Fingerprint  string
	ConnectionID int64
	Otp          string
}

type SignUpDNDParams struct {
	NickName   string `json:"nickname" db:"nickname"`
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	IsDnd      bool   `json:"is_dnd" db:"is_dnd"`
}

type ClientId struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}

type AddTotpSecretParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	SecretTotp string `json:"-" db:"secret_totp"`
}

type SignUpParams struct {
	UA       string `json:"ua" db:"ua"`
	IP       string `json:"ip"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`
}

type ConfirmPhoneParams struct {
	Hash        string `json:"hash" db:"hash"`
	CodeConfirm string `json:"code_confirm" db:"code_confirm"`
}

type ValidateParams struct {
	Signature string
	Message   string
	Public    string
}

type ConfirmMailParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Hash       string `json:"hash" db:"code_need"`
}
type ConfirmTGParams struct {
	Hash string `json:"hash" db:"code_need"`
}

type RequestToConfirmTGParams struct {
	ClientID    int64  `json:"client_id"`
	LanguageIso string `json:"language_iso"`
}

type IsKYCParams struct {
	ClientID int64 `json:"clients_id" db:"clients_id"`
}

type TokenGenerate struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Email      string `json:"email"`
}

type IsPhoneConfirmParams struct {
	ClientID int64 `json:"clients_id" db:"clients_id"`
}

type IsEmailConfirmParams struct {
	ClientID int64 `json:"clients_id" db:"clients_id"`
}

type GetActiveSessionParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}

type ChangePasswordParams struct {
	ClientUuid  string `json:"client_uuid" db:"client_uuid"`
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password,omitempty"`
}
type ChangeEmailParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	NewEmail   string `json:"new_email" db:"email"`
}

type Password struct {
	Pas string `json:"-" db:"password"`
}

type VerCodeParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	CodeInput  string `json:"code_input"`
	Type       string `json:"type"`
}

type SendCodeParams struct {
	ClientUuid  string `json:"client_uuid" db:"client_uuid"`
	LanguageIso string `json:"language_iso"`
	Type        string `json:"type"`
}

type ChangeTgParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	NewTg      string `json:"tg_user_name" db:"tg"`
}

type ChangeNicknameParams struct {
	ClientUuid  string `json:"client_uuid" db:"client_uuid"`
	NewNickname string `json:"new_nickname" db:"nickname"`
	OldNickname string `json:"old_nickname,omitempty" db:"old_nickname"`
	Refresh     string `json:"refresh"`
}

type CompareKeyWordParams struct {
	ClientID int64  `json:"clients_id" db:"clients_id"`
	KeyWord  string `json:"key_word" db:"tg"`
}

type GetEmailParams struct {
	ClientUuid  string `json:"client_uuid" db:"client_uuid"`
	LanguageIso string `json:"language_iso"`
}

type RecoveryInitParams struct {
	ClientUuid  string `json:"client_uuid" db:"client_uuid"`
	LanguageIso string `json:"language_iso,omitempty"`
	Email       string `json:"email"`
}

type RecoveryConfirmParams struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}

type ResponseGetClientPrivateModel struct {
	Tg string `json:"tg"`
}

type GetAuthLevelParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}

type GetPhoneParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}

type AddTotpParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}

type VerifyTotpParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Token      string `json:"totp_token"`
}

type AddTotpRepoParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Secret     string `json:"secret"`
}

type GetClientParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}
type GetClientAuthLevelParams struct {
	ClientID int64 `json:"clients_id" db:"clients_id"`
}
type ChangePhoneParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	NewPhone   string `json:"new_phone" db:"phone"`
}

type RequestToConfirmMailParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	CodeInput  string `json:"code_input"`
}

type RequestToConfirmPhoneParams struct {
	ClientID int64  `json:"clients_id" db:"clients_id"`
	Phone    string `json:"email" db:"email"`
}

type KycConfirmInitParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}

type KycConfirmParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Hash       string `json:"hash"`
}

type HashKycConfirmBody struct {
	ClientID string `json:"clientId"`
	CallBack string `json:"callbackUrl"`
}

type CreateClientUidParamsRepo struct {
	ClientID int64 `json:"client_id" db:"id"`
}
