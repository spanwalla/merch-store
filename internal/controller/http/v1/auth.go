package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/spanwalla/merch-store/internal/service"
	"net/http"
)

type authRoutes struct {
	authService service.Auth
}

type getTokenInput struct {
	Username string `json:"username" validate:"required,min=4,max=64"`
	Password string `json:"password" validate:"required,password"`
}

func newAuthRoutes(g *echo.Group, authService service.Auth) {
	r := &authRoutes{authService}

	g.POST("", r.getToken)
}

func (r *authRoutes) getToken(c echo.Context) error {
	var input getTokenInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	token, err := r.authService.GenerateToken(c.Request().Context(), service.AuthGenerateTokenInput{
		Name:     input.Username,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrWrongPassword) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Token string `json:"token"`
	}

	return c.JSON(http.StatusOK, response{token})
}
