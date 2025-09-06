package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/repositories"
)

type CinemaHandler struct {
	cr *repositories.CinemaRepository
}

func NewCinemaHandler(cr *repositories.CinemaRepository) *CinemaHandler {
	return &CinemaHandler{cr: cr}
}

func newScheduleResponse(res []models.CinemaSchedule, success bool, error string) models.ScheduleResponse {
	return models.ScheduleResponse{Result: res, Success: success, Error: error}
}

func (c *CinemaHandler) HandlerSchedule(ctx *gin.Context) {
	schedule, err := c.cr.GetSchedule(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newScheduleResponse(
			nil, false, "server unable to get cinema schedule",
		))
	}

	if len(schedule) == 0 {
		ctx.JSON(http.StatusNoContent, newScheduleResponse(
			[]models.CinemaSchedule{}, true, "cinema schedules is empty",
		))
		return
	}

	ctx.JSON(http.StatusOK, newScheduleResponse(
		schedule, true, "",
	))
}

func newAvailSeatsRepsonse(res []models.AvailSeat, success bool, err string) models.AvailSeatsResponse {
	return models.AvailSeatsResponse{Result: res, Success: success, Error: err}
}

func (c *CinemaHandler) HandlerSeats(ctx *gin.Context) {
	seats, err := c.cr.GetAvailSeats(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newAvailSeatsRepsonse(
			nil, false, "server unable to get available seats",
		))
		return
	}

	if len(seats) == 0 {
		ctx.JSON(http.StatusNoContent, newAvailSeatsRepsonse(
			[]models.AvailSeat{}, true, "there is available seat",
		))
		return
	}

	ctx.JSON(http.StatusOK, newAvailSeatsRepsonse(
		seats, true, "",
	))
}
