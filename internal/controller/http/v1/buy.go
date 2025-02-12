package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/spanwalla/merch-store/internal/service"
	"net/http"
)

type buyRoutes struct {
	paymentService service.Payment
}

type buyItemInput struct {
	Item string `param:"item" validate:"required,max=16"`
}

func newBuyRoutes(g *echo.Group, paymentService service.Payment) {
	r := &buyRoutes{paymentService}

	g.GET("/:item", r.buyItem)
}

func (r *buyRoutes) buyItem(c echo.Context) error {
	var input buyItemInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid params")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.paymentService.BuyItem(c.Request().Context(), service.PaymentBuyItemInput{
		UserId:   c.Get(userIdCtx).(int),
		ItemName: input.Item,
	})
	if err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, service.ErrNotEnoughBalance) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return err
	}

	return c.NoContent(http.StatusOK)
}
