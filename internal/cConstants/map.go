package cConstants

var CodesConfirmExpiredTime = map[string]string{
	"recovery_by_email": "1 hour",
	"email_confirm":     "1 hour",
	"KYC":               "1 week",
	"phone_confirm":     "1 hour",
	"confirm_withdraw":  "10 minute",
}

const (
	Email           = "email_confirm"
	KYC             = "KYC"
	RecoveryByEmail = "recovery_by_email"
)
