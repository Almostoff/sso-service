package usecase

import (
	"AuthService/internal/cConstants"
	"AuthService/internal/iConnection"
	"AuthService/internal/model"
	"AuthService/pkg/secure"
	"errors"
	"fmt"
	"log"
	"strings"
)

type iConnectionUseCase struct {
	repo iConnection.Repository
}

func NewIConnectionUsecase(repo iConnection.Repository) iConnection.UseCase {
	return &iConnectionUseCase{repo: repo}
}

func (u *iConnectionUseCase) AddRedirectUrl(params *model.Response) {
	//TODO implement me
	panic("implement me")
}

func (u *iConnectionUseCase) Validate(params *iConnection.ValidateParams) (*bool, error) {

	connection, err := u.repo.GetServiceByPublic(&iConnection.GetServiceByPublicParams{Public: params.Public})
	if err != nil {
		return &cConstants.False, err
	}

	body := createRequestBody(params.Timestamp, strings.TrimRight(params.Message, "\n"))
	hash := secure.CalcSignature(*connection.Private, body)
	if hash != params.Signature {
		return &cConstants.False, fmt.Errorf("hash {%s} != signature {%s}", hash, params.Signature)
	}

	return &cConstants.True, nil
}

func (u *iConnectionUseCase) GetInnerConnection(params *iConnection.GetInnerConnectionParams) (*iConnection.InnerConnection, error) {

	connection, err := u.repo.GetInnerConnection(&iConnection.GetInnerConnectionParams{Name: params.Name})
	if err != nil {
		log.Println("here yeah?")
		return &iConnection.InnerConnection{}, err
	}
	if connection == nil {
		return &iConnection.InnerConnection{}, errors.New("connection is nil somehow")
	}

	return connection, nil
}

func createRequestBody(timestamp, jsonBody string) string {
	return timestamp + jsonBody
}
