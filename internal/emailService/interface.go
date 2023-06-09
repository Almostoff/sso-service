package emailService

import (
	"AuthService/internal/model"
)

type Email interface {
	SendMailConfirm(params *SendMailConfirmParams) error
	SendMail(params any) error
	SendPhoneConfirm(params *SendPhoneConfirmParams) *ResponseSendPhoneConfirm
	SendRecoveryLink(params *SendRecoveryLinkParams) *model.CodeModel
}
