package http

import (
	"AuthService/config"
	"AuthService/internal/SSO"
	"AuthService/internal/model"
	"AuthService/pkg/logger"
	"AuthService/pkg/loggerService"
	"AuthService/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type ClientHandler struct {
	logger   *logger.ApiLogger
	cfg      *config.Config
	clientUC SSO.UseCase
}

func NewClientsHandlers(cfg *config.Config, logger *logger.ApiLogger, clientUC SSO.UseCase) *ClientHandler {
	return &ClientHandler{cfg: cfg, logger: logger, clientUC: clientUC}
}

func (c ClientHandler) SignUp() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.SignUpParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.SignUp(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) ClientSignIn() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.SignInParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.SignIn(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) ClientSignInTg() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.SignInTGParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.SignInTg(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

//func (c ClientHandler) ConfirmPhone() fiber.Handler {
//	return func(ctx *fiber.Ctx) error {
//		var params model.ConfirmPhoneParams
//		if err := utils.ReadRequest(ctx, &params); err != nil {
//			return ctx.SendStatus(fiber.StatusBadRequest)
//		}
//		data := c.clientUC.ConfirmPhone(&params)
//		data.Error = model.GetResponseCode(data.Error, data.Data, params)
//		return ctx.JSON(data)
//	}
//}

//func (c ClientHandler) RequestToConfirmPhone() fiber.Handler {
//	return func(ctx *fiber.Ctx) error {
//		var params model.RequestToConfirmPhoneParams
//		if err := utils.ReadRequest(ctx, &params); err != nil {
//			return ctx.SendStatus(fiber.StatusBadRequest)
//		}
//		data := c.clientUC.RequestToConfirmPhone(&params)
//		data.Error = model.GetResponseCode(data.Error, data.Data, params)
//		return ctx.JSON(data)
//	}
//}
//
//func (c ClientHandler) CompareKeyWord() fiber.Handler {
//	return func(ctx *fiber.Ctx) error {
//		var params model.CompareKeyWordParams
//		if err := utils.ReadRequest(ctx, &params); err != nil {
//			return ctx.SendStatus(fiber.StatusBadRequest)
//		}
//		data := c.clientUC.CompareKeyWord(&params)
//		data.Error = model.GetResponseCode(data.Error, data.Data, params)
//		return ctx.Status(handleStatus(data.Error)).JSON(data)
//	}
//}

func (c ClientHandler) ConfirmMail() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.ConfirmMailParams
		params.Hash = ctx.Params("hash")
		data := c.clientUC.ConfirmMail(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) ConfirmKyc() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.KycConfirmParams
		var body model.KYCResponse
		if err := utils.ReadRequest(ctx, &body); err != nil {
			badBody := string(ctx.Body())
			go loggerService.GetInstance().DevLog(fmt.Sprintf("!🙀! CallBack от idefi BADREQUEST:\n %+v", badBody), 8)
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		body.FileUrls.BACK = "****"
		body.FileUrls.FRONT = "****"
		body.FileUrls.FACE = "***"
		go loggerService.GetInstance().DevLog(fmt.Sprintf("! 🙀! CallBack от idefi: %+v", body), 8)
		if body.Status.Overall != "APPROVED" {
			go loggerService.GetInstance().DevLog(fmt.Sprintf("!🙀! CallBack от idefi not APPROVED: %+v", body), 8)
			return ctx.SendStatus(200)
		}
		params.Hash = body.ClientId
		data := c.clientUC.KycConfirm(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)

	}
}

func (c ClientHandler) RequestToConfirmMail() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.RequestToConfirmMailParams
		params.CodeInput = ctx.Get("hash")
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.RequestToConfirmMail(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) ChangePassword() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.ChangePasswordParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if params.ClientUuid == "" {
			uuid := c.getUserUuid(ctx)
			if uuid == "" {
				return ctx.SendStatus(fiber.StatusUnauthorized)
			}
			params.ClientUuid = uuid
		}
		data := c.clientUC.ChangePasswordWithOldCheck(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) ChangeEmail() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.ConfirmMailParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.ConfirmMail(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) ChangePhone() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.ChangePhoneParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		uuid := c.getUserUuid(ctx)
		if uuid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		params.ClientUuid = uuid
		data := c.clientUC.ChangePhone(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) ChangeTg() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.ChangeTgParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		uuid := c.getUserUuid(ctx)
		if uuid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		params.ClientUuid = uuid
		data := c.clientUC.ChangeTg(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) ChangeNickname() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.ChangeNicknameParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		uuid := c.getUserUuid(ctx)
		if uuid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		params.ClientUuid = uuid
		data := c.clientUC.ChangeNickname(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) RecoveryInit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.RecoveryInitParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.RecoveryInit(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) RecoveryConfirm() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.RecoveryConfirmParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.RecoveryConfirm(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) AddTotp() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.AddTotpParams
		uuid := c.getUserUuid(ctx)
		if uuid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		params.ClientUuid = uuid
		data := c.clientUC.AddTotp(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) VerifyTotp() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.VerifyTotpParams
		uuid := c.getUserUuid(ctx)
		if uuid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		params.ClientUuid = uuid
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.VerifyTotp(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) VerifyTotpInit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.VerifyTotpParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.VerifyTotpInit(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) ConfirmKycInit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.KycConfirmInitParams
		uuid := c.getUserUuid(ctx)
		if uuid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		params.ClientUuid = uuid
		data := c.clientUC.KycConfirmInit(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) IsCodeValid() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.VerCodeParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.IsCodeValid(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) AddCode() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.WriteAuthCodeParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.AddCode(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) GetAuthLevel() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		uuid := c.getUserUuid(ctx)
		if uuid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		data := c.clientUC.GetAuthLevel(&uuid)
		data.Error = model.GetResponseCode(data.Error, data.Data, uuid)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) GetClient() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.GetClientParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.GetClientPrivateInfo(params.ClientUuid)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) ValidateAccess() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.ValidateAccessTokenParams
		//uuid := c.getUserUuid(ctx)
		//if uuid == "" {
		//	return ctx.SendStatus(fiber.StatusUnauthorized)
		//}
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.ValidateAccessToken(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) ValidateRefresh() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.ValidateRefreshTokenParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.ValidateRefreshToken(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) RefreshAccess() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.RefreshAccessTokenParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := c.clientUC.RefreshAccessToken(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.JSON(data)
	}
}

func (c ClientHandler) Logout() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.LogoutParams
		params.UA = ctx.Get("ua")
		params.Access = ctx.Get("access")
		if params.Access == "" {
			if err := utils.ReadRequest(ctx, &params); err != nil {
				return ctx.SendStatus(fiber.StatusBadRequest)
			}
		}
		data := c.clientUC.Logout(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) GetActiveSessions() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.GetActiveSessionParams
		uuid := c.getUserUuid(ctx)
		if uuid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		params.ClientUuid = uuid
		data := c.clientUC.GetActiveSession(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func (c ClientHandler) DeleteSession() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.GetActiveSessionParams
		uuid := c.getUserUuid(ctx)
		if uuid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		params.ClientUuid = uuid
		data := c.clientUC.GetActiveSession(&params)
		data.Error = model.GetResponseCode(data.Error, data.Data, params)
		return ctx.Status(handleStatus(data.Error)).JSON(data)
	}
}

func handleStatus(cErr *model.CodeModel) int {
	if cErr.IsError {
		return int(cErr.StandardCode)
	} else {
		return 200
	}
}

func (c ClientHandler) getUserUuid(ctx *fiber.Ctx) string {
	access := ctx.Get("access")
	if access == "" {
		return ""
	}
	email := c.clientUC.DecodeToken(access)
	if email == "" {
		return ""
	}
	uuid, cErr := c.clientUC.GetClientUuidByEmail(email)
	if cErr.IsError {
		return ""
	}
	return *uuid
}
