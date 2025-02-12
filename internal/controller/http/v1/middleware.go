package v1

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spanwalla/merch-store/internal/service"
	"net/http"
	"strings"
)

const userIdCtx = "userId"

type AuthMiddleware struct {
	authService service.Auth
}

func (h *AuthMiddleware) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c.Request())
		if !ok {
			log.Errorf("AuthMiddleware.UserIdentity - bearerToken: %v", ErrInvalidAuthHeader)
			newErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			return nil
		}

		userId, err := h.authService.VerifyToken(token)
		if err != nil {
			log.Errorf("AuthMiddleware.UserIdentity - ParseToken: %v", err)
			newErrorResponse(c, http.StatusUnauthorized, ErrCannotParseToken.Error())
			return err
		}

		c.Set(userIdCtx, userId)

		return next(c)
	}
}

func bearerToken(req *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := req.Header.Get(echo.HeaderAuthorization)
	if len(header) == 0 {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
