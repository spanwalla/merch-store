package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/spanwalla/merch-store/internal/service"
	"net/http"
)

type infoRoutes struct {
	userReportService service.UserReport
}

func newInfoRoutes(g *echo.Group, userReportService service.UserReport) {
	r := &infoRoutes{userReportService}

	g.GET("", r.getReport)
}

func (r *infoRoutes) getReport(c echo.Context) error {
	report, err := r.userReportService.Get(c.Request().Context(), c.Get(userIdCtx).(int))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, report)
}
