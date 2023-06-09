package httpServer

import (
	clientHttp "AuthService/internal/SSO/delivery/http"
	clientsRepository "AuthService/internal/SSO/repository"
	clientsUseCase "AuthService/internal/SSO/usecase"
	"AuthService/internal/cConstants"
	"AuthService/internal/emailService"
	"AuthService/internal/iConnection"
	iConnectionRepository "AuthService/internal/iConnection/repository"
	iConnectionUC "AuthService/internal/iConnection/usecase"
	"AuthService/internal/middleware"
	redisRepo "AuthService/internal/redisUsers/repository"
	redisUC "AuthService/internal/redisUsers/usecase"
	"AuthService/pkg/logger"
	"AuthService/pkg/storage"
	"AuthService/pkg/storage/redis"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	serverLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (s *Server) MapHandlers(app *fiber.App, logger *logger.ApiLogger) error {
	rdb := redis.InitRedisClient(s.cfg)
	db, err := storage.InitPsqlDB(s.cfg)
	if err != nil {
		return err
	}

	iConnectionRepo := iConnectionRepository.NewPostgresRepository(db, s.shield)
	iConnectionUC := iConnectionUC.NewIConnectionUsecase(iConnectionRepo)
	textConnection, err := iConnectionUC.GetInnerConnection(&iConnection.GetInnerConnectionParams{Name: &cConstants.EmailService})
	if err != nil {
		log.Println("here?")
		return err
	}
	textSer := emailService.UsersTextsClient(&emailService.GetClientParams{
		Private: textConnection.Private,
		BaseUrl: textConnection.BaseUrl,
		Public:  textConnection.Public,
		Config:  s.cfg,
	})

	redisRepo := redisRepo.NewRedisRepository(rdb, s.shield)
	redisUC := redisUC.RedisUseCase(redisRepo, s.shield)

	usersRepo := clientsRepository.NewPostgresRepository(db, s.shield)
	usersUC := clientsUseCase.NewUsersUsecase(logger, usersRepo, s.shield, redisUC, textSer, s.cfg.Kyc, s.cfg.Server)
	usersHandlers := clientHttp.NewClientsHandlers(s.cfg, logger, usersUC)

	app.Use(serverLogger.New())
	if _, ok := os.LookupEnv("LOCAL"); !ok {
		app.Use(recover.New())
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	mw := middleware.NewMDWManager(s.cfg, usersUC, iConnectionUC, logger)
	clientHttp.MapAdminsRoutes(app, usersHandlers, mw)

	return nil
}
