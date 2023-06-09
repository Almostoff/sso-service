package emailService

import (
	"AuthService/config"
	"AuthService/internal/cConstants"
	"AuthService/internal/model"
	"AuthService/pkg/utils"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	cfg        *config.Config
	httpClient *resty.Client
}

type GetClientParams struct {
	Config  *config.Config
	Public  *string
	Private *string
	BaseUrl *string
}

func UsersTextsClient(params *GetClientParams) Email {
	return &Client{
		cfg: params.Config,
		httpClient: resty.New().OnBeforeRequest(utils.SignatureMiddleware((*utils.GetClientParams)(params))).
			//EnableTrace().SetDebug(true).SetBaseURL(*params.BaseUrl),
			EnableTrace().SetDebug(true).SetBaseURL("http://localhost:8080"),
	}
}

func (c Client) SendMail(params any) error {
	var responseModel ResponseSendMailConfirm
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(&params).Post(cConstants.TextSendMail)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if response == nil {
		return errors.New("response is nil somehow")
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return fmt.Errorf("status code {%d}", statusCode)
	}

	return nil
}

func (c Client) SendMailConfirm(params *SendMailConfirmParams) error {
	var responseModel ResponseSendMailConfirm
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.TextSendMail)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if response == nil {
		return errors.New("response is nil somehow")
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return fmt.Errorf("status code {%d}", statusCode)
	}

	return nil
}

func (c Client) SendRecoveryLink(params *SendRecoveryLinkParams) *model.CodeModel {
	var responseModel model.CodeModel
	response, err := c.httpClient.R().SetBody(params).Post(cConstants.TextSendMail)
	if err != nil {
		fmt.Println(err.Error())
		responseModel.StandardMessage = err.Error()
		responseModel.StandardCode = cConstants.StatusInternalServerError
		responseModel.InternalCode = cConstants.StatusInternalServerError
		return &responseModel
	}
	if response == nil {
		responseModel.StandardMessage = "nil response"
		responseModel.StandardCode = cConstants.StatusInternalServerError
		responseModel.InternalCode = cConstants.StatusInternalServerError
		return &responseModel
	}
	statusCode := int64(response.StatusCode())
	if statusCode > 399 {
		responseModel.StandardMessage = fmt.Sprintf("respose status code %d", statusCode)
		responseModel.StandardCode = cConstants.StatusInternalServerError
		responseModel.InternalCode = cConstants.StatusInternalServerError
		return &responseModel
	}

	return &responseModel
}

func (c Client) SendPhoneConfirm(params *SendPhoneConfirmParams) *ResponseSendPhoneConfirm {
	//TODO implement me
	panic("implement me")
}
