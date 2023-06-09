package model

type ChangeNicknameData struct {
	Access string `json:"Access"`
}

type SignUpData struct {
	AccessToken  string `json:"Access"`
	RefreshToken string `json:"Refresh"`
}

type SignInData struct {
	AccessToken  string `json:"Access"`
	RefreshToken string `json:"Refresh"`
}

type SignInWith2faData struct {
	AccessToken string `json:"Access"`
}

type SuccessData struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type SuccessKycData struct {
	Success bool   `json:"success"`
	Hash    string `json:"hash"`
	Message string `json:"message,omitempty"`
}

type GetEmailData struct {
	Email string `json:"email" db:"email"`
}
type AuthLevelData struct {
	AuthLevel *AuthLevel `json:"auth_level"`
}

type GetPhoneData struct {
	Phone string `json:"phone" db:"phone"`
}

type TotpModel struct {
	File        []byte `json:"file"`
	AccountName string `json:"account_name"`
	Secret      string `json:"secret"`
	Link        string `json:"link"`
}

type ResponseSuccessModel struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
