package model

type Response struct {
	Error   *CodeModel   `json:"error"`
	Success *Success     `json:"success"`
	Data    *interface{} `json:"data,omitempty"`
}

type ResponseSignIn struct {
	Error   *CodeModel  `json:"error"`
	Success *Success    `json:"success"`
	Data    *SignInData `json:"data"`
}

type ResponseSignUp struct {
	Error   *CodeModel  `json:"error"`
	Success *Success    `json:"success"`
	Data    *SignUpData `json:"data"`
}

type ResponseAuthLevel struct {
	Error   *CodeModel `json:"error"`
	Success *Success   `json:"success"`
	Data    *AuthLevel `json:"data"`
}

type ResponseClient struct {
	Error   *CodeModel `json:"error"`
	Success *Success   `json:"success"`
	Data    *Client    `json:"data"`
}

type ResponseAddTotp struct {
	Error   *CodeModel `json:"error"`
	Success *Success   `json:"success"`
	Data    *TotpModel `json:"data"`
}

type ResponseGetActiveSession struct {
	Error   *CodeModel `json:"error"`
	Success *Success   `json:"success"`
	Data    *[]Session `json:"data"`
}

type ResponseValidateToken struct {
	Error *CodeModel            `json:"error"`
	Data  *ResponseSuccessModel `json:"data"`
}

type ResponseRefreshAccessToken struct {
	Error *CodeModel `json:"error"`
	Data  *Access    `json:"data"`
}
