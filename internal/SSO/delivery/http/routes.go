package http

import (
	"AuthService/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func MapAdminsRoutes(router fiber.Router, c *ClientHandler, mw *middleware.MDWManager) {
	clientRouter := router.Group("/sso")
	clientRouter.Post("/sign_up", c.SignUp())
	clientRouter.Post("/sign_in", c.ClientSignIn())

	privateRouter := clientRouter.Group("/my")
	//mw.ValidateAccess())
	privateRouter.Post("/logout", c.Logout())

	getgr := privateRouter.Group("/get")
	getgr.Get("/private_info", c.GetClient())
	getgr.Get("/verify_level", c.GetAuthLevel())

	chgr := privateRouter.Group("/change")
	chgr.Post("/phone", c.ChangePhone())
	chgr.Post("/password", c.ChangePassword())
	chgr.Post("/tg", c.ChangeTg())
	chgr.Post("/nickname", c.ChangeNickname())

	recgr := clientRouter.Group("/recovery")
	recgr.Post("/init_by_email", c.RecoveryInit())
	recgr.Post("/by_email", c.RecoveryConfirm())

	confgr := privateRouter.Group("/confirm")
	confgradd := confgr.Group("/add")
	confgradd.Get("/totp", c.AddTotp())
	confgradd.Get("/email", c.RequestToConfirmMail())
	confgradd.Post("/code", c.AddCode())
	confgradd.Get("/kyc", c.ConfirmKycInit())
	confgradd.Post("/key_word", c.AddCode())
	confgrver := confgr.Group("/verify")
	confgrver.Post("/totp", c.VerifyTotp())
	confgrver.Post("/totp_init", c.VerifyTotpInit())
	confgrver.Post("/code", c.IsCodeValid())
	confgrver.Get("/email/:hash", c.ConfirmMail())
	confgrver.Post("/kyc", c.ConfirmKyc())

	tkgr := clientRouter.Group("/token")
	valgr := tkgr.Group("/validate")
	valgr.Post("/access", c.ValidateAccess())
	valgr.Post("/refresh", c.ValidateRefresh())
	rfgr := tkgr.Group("/update")
	rfgr.Post("/access", c.RefreshAccess())

	sesgr := privateRouter.Group("/session")
	sesgr.Delete("/delete/all_except_current", c.GetActiveSessions())
	sesgr.Delete("/delete/:id", c.DeleteSession())
	sesgr.Get("/", c.GetActiveSessions())

	serviceRouter := clientRouter.Group("/service")
	//mw.VerifySignatureMiddleware())
	serviceRouter.Post("/get_client", c.GetClient())
}
