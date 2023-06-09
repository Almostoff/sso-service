package tg

type Email interface {
	SendMailConfirm(params *SendMailConfirmParams) error
}
