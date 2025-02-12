package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/spanwalla/merch-store/internal/service"
	"net/http"
)

type sendRoutes struct {
	paymentService service.Payment
}

type sendCoinInput struct {
	ToUser string `json:"toUser" validate:"required,min=4,max=64"`
	Amount int    `json:"amount" validate:"required,gt=0"`
}

func newSendRoutes(g *echo.Group, paymentService service.Payment) {
	r := &sendRoutes{paymentService}

	g.POST("", r.sendCoin)
}

func (r *sendRoutes) sendCoin(c echo.Context) error {
	var input sendCoinInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.paymentService.Transfer(c.Request().Context(), service.PaymentTransferInput{
		FromUserId: c.Get(userIdCtx).(int),
		ToUserName: input.ToUser,
		Amount:     input.Amount,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, service.ErrNotEnoughBalance) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, service.ErrSelfTransfer) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return err
	}

	return c.NoContent(http.StatusOK)
}
