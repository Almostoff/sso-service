package middleware

import (
	"AuthService/config"
	"AuthService/internal/SSO"
	"AuthService/internal/iConnection"
	"AuthService/pkg/logger"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type MDWManager struct {
	cfg           *config.Config
	usersUC       SSO.UseCase
	iConnectionUC iConnection.UseCase
	logger        *logger.ApiLogger
}

func NewMDWManager(cfg *config.Config, usersUC SSO.UseCase, iConnectionUC iConnection.UseCase, logger *logger.ApiLogger) *MDWManager {
	return &MDWManager{usersUC: usersUC, cfg: cfg, iConnectionUC: iConnectionUC, logger: logger}
}

func (mw *MDWManager) VerifySignatureMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var (
			publicKey string = c.Get("ApiPublic")
			signature string = c.Get("Signature")
			timestamp string = c.Get("Timestamp")
			method    string = c.Method()
			err       error
			message   string
		)

		if method == "POST" {
			message = string(c.Body())
			message = strings.ReplaceAll(message, "\n", "")
		} else {
			message = publicKey
		}
		// ------------------------------------------------------------------------------------------------

		if timestamp == "" {
			err = errors.New("timestamp is required")
			mw.logger.Error(err)
			return err
		}

		isValid, err := mw.iConnectionUC.Validate(&iConnection.ValidateParams{
			Signature: signature,
			Public:    publicKey,
			Message:   message,
			Timestamp: timestamp,
		})
		if err != nil || !*isValid {
			mw.logger.Error(err)
			return err
		}
		mw.logger.Info("OK")
		return c.Next()
	}
}
