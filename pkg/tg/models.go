package tg

type Params struct {
	ApiPublic string
	Signature string
	TimeStamp string
}

type ConnData struct {
	ApiPublic string
	ApiSecret string
	TimeStamp string
}

type SendMailConfirmParams struct {
	TypeEmail   string `json:"type_email" bson:"type_email"`
	Email       string `db:"email" json:"email_receiver"`
	LanguageIso string `json:"language_iso"`
	Link        string `json:"link"`
}

type SendRecoveryLinkParams struct {
	TypeEmail   string `json:"type_email" bson:"type_email"`
	Email       string `db:"email" json:"email_receiver"`
	LanguageIso string `json:"language_iso"`
	Link        string `json:"link"`
}

type SendPhoneConfirmParams struct {
	Phone string `db:"phone" json:"phone"`
}

type ResponseSendMailConfirm struct {
	Error *model.CodeModel `json:"error"`
	Data  *interface{}     `json:"data"`
}

type ResponseSendPhoneConfirm struct {
	Error *model.CodeModel `json:"error"`
	Data  *interface{}     `json:"data"`
}
