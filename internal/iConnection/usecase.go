package iConnection

import "AuthService/internal/model"

type UseCase interface {
	GetInnerConnection(params *GetInnerConnectionParams) (*InnerConnection, error)
	Validate(params *ValidateParams) (*bool, error)
	AddRedirectUrl(params *model.Response)
}
