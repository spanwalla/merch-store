package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spanwalla/merch-store/internal/service"
	"net/http"
	"os"
)

func ConfigureRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.CORS())
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	authGroup := handler.Group("/api/auth")
	{
		newAuthRoutes(authGroup, services.Auth)
	}

	authMiddleware := &AuthMiddleware{services.Auth}
	protectedGroup := handler.Group("/api", authMiddleware.UserIdentity)
	{
		newInfoRoutes(protectedGroup.Group("/info"), services.UserReport)
		newBuyRoutes(protectedGroup.Group("/buy"), services.Payment)
		newSendRoutes(protectedGroup.Group("/sendCoin"), services.Payment)
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("v1 - setLogsFile - os.OpenFile: %v", err)
	}
	return file
}
