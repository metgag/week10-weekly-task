package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/metgag/koda-weekly10/internals/utils"
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

// HandleCinemaSchedule godoc
//
//	@Summary		Get cinema schedules
//	@Description	Retrieve all available cinema schedules from the server
//	@Tags			cinemas
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.ScheduleResponse	"Schedules fetched successfully or no schedules available"
//	@Failure		500	{object}	models.ScheduleResponse	"Internal server error while fetching the schedule"
//	@Security		BearerAuth
//	@Router			/cinemas/schedules [get]
func (c *CinemaHandler) HandlerSchedule(ctx *gin.Context) {
	schedule, err := c.cr.GetSchedule(ctx.Request.Context())
	if err != nil {
		utils.PrintError("CINEMA SCHEDULE SERVER ERROR", 8, err)
		ctx.JSON(http.StatusInternalServerError, newScheduleResponse(
			nil, false, "server unable to get cinema schedule",
		))
	}

	if len(schedule) == 0 {
		utils.PrintError("NO CINEMA SCHEDULES", 12, nil)
		ctx.JSON(http.StatusOK, newScheduleResponse(
			[]models.CinemaSchedule{}, true, "no cinema schedules available",
		))
		return
	}

	ctx.JSON(http.StatusOK, newScheduleResponse(
		schedule, true, "",
	))
}

func newAvailSeatsRepsonse(res []models.Seat, success bool, err string) models.AvailSeatsResponse {
	return models.AvailSeatsResponse{Result: res, Success: success, Error: err}
}

// HandleCinemaSeats godoc
//
//	@Summary		Get available seats
//	@Description	Retrieve all available seats for a specific cinema schedule
//	@Tags			cinemas
//	@Accept			json
//	@Produce		json
//	@Param			schedule_id	path		int							true	"The ID of the cinema schedule"
//	@Success		200			{object}	models.AvailSeatsResponse	"Available seats retrieved successfully or no seats available"
//	@Failure		400			{object}	models.AvailSeatsResponse	"Invalid schedule ID format"
//	@Failure		500			{object}	models.AvailSeatsResponse	"Internal server error while fetching available seats"
//	@Security		BearerAuth
//	@Router			/cinemas/{schedule_id}/seats [get]
func (c *CinemaHandler) HandlerSeats(ctx *gin.Context) {
	scheduleIdStr := ctx.Param("schedule_id")
	scheduleId, err := strconv.Atoi(scheduleIdStr)
	if err != nil {
		utils.PrintError("INVALID SCHEDULE ID", 8, err)
		ctx.JSON(http.StatusBadRequest, newAvailSeatsRepsonse(
			nil, false, "invalid schedule ID format",
		))
		return
	}

	seats, err := c.cr.GetAvailSeats(ctx, scheduleId)
	if err != nil {
		utils.PrintError("CINEMA AVAIL SEATS SERVER ERROR", 8, err)
		ctx.JSON(http.StatusInternalServerError, newAvailSeatsRepsonse(
			nil, false, "server unable to get available seats",
		))
		return
	}

	if len(seats) == 0 {
		ctx.JSON(http.StatusOK, newAvailSeatsRepsonse(
			[]models.Seat{}, true, "no available seats for this schedule",
		))
		return
	}

	ctx.JSON(http.StatusOK, newAvailSeatsRepsonse(
		seats, true, "",
	))
}

func newCinemaAndTimeResponse(res models.CinemaAndTime, success bool, err string) models.CinemaAndTimeResponse {
	return models.CinemaAndTimeResponse{Result: res, Success: success, Error: err}
}

func (c *CinemaHandler) HandlerCinemaNameAndTime(ctx *gin.Context) {
	scheduleIdStr := ctx.Param("schedule_id")
	scheduleId, err := strconv.Atoi(scheduleIdStr)
	if err != nil {
		utils.PrintError("INVALID SCHEDULE ID", 8, err)
		ctx.JSON(http.StatusBadRequest, newCinemaAndTimeResponse(
			models.CinemaAndTime{}, false, "invalid schedule ID format",
		))
		return
	}

	result, err := c.cr.GetCinemaNameAndTime(ctx.Request.Context(), scheduleId)
	if err != nil {
		utils.PrintError("ERROR GET CINEMA AND TIME", 16, err)
		ctx.JSON(http.StatusInternalServerError, newCinemaAndTimeResponse(
			models.CinemaAndTime{}, false, "SERVER ERROR GET CINEMA AND TIME",
		))
		return
	}

	ctx.JSON(http.StatusOK, newCinemaAndTimeResponse(
		result, true, "",
	))
}
