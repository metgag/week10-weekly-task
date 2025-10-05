package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/metgag/koda-weekly10/internals/utils"
)

type MovieHandler struct {
	mr *repositories.MovieRepository
}

func NewMovieHandler(mr *repositories.MovieRepository) *MovieHandler {
	return &MovieHandler{mr: mr}
}

func newMoviesResponse(result []models.MovieFilter, success bool, error string) models.MoviesResponse {
	return models.MoviesResponse{Result: result, Success: success, Error: error}
}

// HandleGetUpcomingMovies godoc
//
//	@Summary		Get upcoming movies
//	@Description	Get list of movies with release date later than the current date
//	@Tags			movies
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.MovieResponse	"Upcoming movies fetched successfully"
//	@Failure		404	{object}	models.MovieResponse	"No upcoming movies found"
//	@Failure		500	{object}	models.MovieResponse	"Internal server error"
//	@Router			/movies/upcoming [get]
func (m *MovieHandler) GetUpcomingMovies(ctx *gin.Context) {
	movies, err := m.mr.GetUpcomingMovies(ctx.Request.Context())
	if err != nil {
		utils.PrintError("UPCOMING MOVIES SERVER ERROR", 8, err)
		ctx.JSON(http.StatusInternalServerError, newMoviesResponse(
			nil, false, "server unable to get movies",
		))
		return
	}

	if len(movies) == 0 {
		utils.PrintError("UPCOMING MOVIES DATA IS EMPTY", 8, nil)
		ctx.JSON(http.StatusNotFound, newMoviesResponse(
			nil, false, "movies data is empty",
		))
		return
	}

	ctx.JSON(http.StatusOK, newMoviesResponse(
		movies, true, "",
	))
}

// HandleGetPopularMovies godoc
//
//	@Summary		get popular movies handler func
//	@Description	get list of movies where popularity higher than 40
//	@Tags			movies
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.MovieResponse
//	@Failure		500	{object}	models.MovieResponse	"internal server error"
//	@Router			/movies/popular [get]
func (m *MovieHandler) GetPopularMovies(ctx *gin.Context) {
	movies, err := m.mr.GetPopularMovies(ctx.Request.Context())
	if err != nil {
		utils.PrintError("POPULAR MOVIES SERVER ERROR", 8, err)
		ctx.JSON(http.StatusInternalServerError, newMoviesResponse(
			nil, false, "server unable to get movies",
		))
		return
	}

	if len(movies) == 0 {
		utils.PrintError("POPULAR MOVIES DATA IS EMPTY", 8, nil)
		ctx.JSON(http.StatusNotFound, newMoviesResponse(
			nil, false, "movies data is empty",
		))
		return
	}

	ctx.JSON(http.StatusOK, newMoviesResponse(
		movies, true, "",
	))
}

func newMovieResponse(res models.Movie, success bool, err string) models.MovieResponse {
	return models.MovieResponse{Result: res, Success: success, Error: err}
}

// HandleGetMovieDetail godoc
//
//	@Summary		get movie detail handler func
//
//	@Description	get movie detail based from id
//
//	@Tags			movies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"movie ID"
//	@Success		200	{object}	models.MovieResponse	"movie detail data"
//	@Failure		400	{object}	models.MovieResponse	"invalid movie id"
//	@Failure		404	{object}	models.MovieResponse	"movie not found"
//	@Failure		500	{object}	models.MovieResponse	"internal server error"
//	@Router			/movies/{id} [get]
func (m *MovieHandler) GetMovieDetail(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.LogCtxError(
			ctx,
			"INVALID MOVIE ID",
			"Internal server error",
			err,
			http.StatusInternalServerError,
		)
		return
	}

	detail, err := m.mr.GetMovieDetail(ctx.Request.Context(), idParam)
	if err != nil {
		utils.LogCtxError(
			ctx,
			"UNABLE GET MOVIE DETAIL",
			"Internal server error",
			err,
			http.StatusInternalServerError,
		)
		return
	}

	ctx.JSON(http.StatusOK, models.NewFullfilledResponse(
		http.StatusOK,
		detail,
	))
}

// HandleGetAllMovies godoc
//
//	@Summary		get all movies (admin)
//	@Description	complete list of movies from the database
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.MovieResponse	"list of all movies"
//	@Failure		500	{object}	models.MovieResponse	"internal server error"
//	@Security		BearerAuth
//	@Router			/admin/movies/ [get]
func (m *MovieHandler) HandleGetAllMovie(ctx *gin.Context) {
	movies, err := m.mr.GetAllMovies(ctx.Request.Context(), "")
	if err != nil {
		utils.PrintError("UNABLE GET MOVIE LIST", 6, err)
		ctx.JSON(http.StatusInternalServerError, newMoviesResponse(
			nil, false, "server unable to get list of movies",
		))
		return
	}
	if len(movies) == 0 {
		utils.PrintError("NO MOVIES ON THE SERVER", 4, nil)
		ctx.JSON(http.StatusInternalServerError, newMoviesResponse(
			nil, false, "no movies on the server",
		))
		return
	}

	ctx.JSON(http.StatusOK, newMoviesResponse(
		movies, true, "",
	))
}

func newDeleteMovieResponse(success bool, res, err string) models.DeleteMovieResponse {
	return models.DeleteMovieResponse{Success: success, Result: res, Error: err}
}

// HandleDeleteMovie godoc
//
//	@Summary		delete a movie w/ ID
//	@Description	soft delete a movie from the database using id
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int							true	"movie ID"
//	@Success		200	{object}	models.DeleteMovieResponse	"Movie deleted successfully"
//	@Failure		400	{object}	models.DeleteMovieResponse	"Invalid movie ID format"
//	@Failure		404	{object}	models.DeleteMovieResponse	"Movie not found"
//	@Failure		500	{object}	models.DeleteMovieResponse	"Server error while deleting movie"
//	@Security		BearerAuth
//	@Router			/admin/movies/{id} [delete]
func (m *MovieHandler) HandleDeleteMovie(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.PrintError("INVALID MOVIE ID", 12, err)
		ctx.JSON(http.StatusBadRequest, newDeleteMovieResponse(
			false, "", "invalid movie ID",
		))
		return
	}

	ctag, err := m.mr.SoftDeleteMovie(ctx.Request.Context(), idParam)
	if err != nil {
		utils.PrintError("UNABLE DELETE MOVIE", 12, err)
		ctx.JSON(http.StatusInternalServerError, newDeleteMovieResponse(
			false, "", "server unable to delete movie",
		))
		return
	}

	if ctag.RowsAffected() == 0 {
		utils.PrintError(fmt.Sprintf("UNABLE TO DELETE MOVIE, NO MOVIE W/ ID %d", idParam), 6, nil)
		ctx.JSON(http.StatusBadRequest, newDeleteMovieResponse(
			false, "", fmt.Sprintf("no movie w/ ID %d", idParam),
		))
		return
	}

	ctx.JSON(http.StatusOK, newDeleteMovieResponse(
		true, fmt.Sprintf("movie w/ ID %d deleted succesfully", idParam), "",
	))
}

func newUpdateMovieResponse(success bool, res, err string) models.DeleteMovieResponse {
	return models.DeleteMovieResponse{Success: success, Result: res, Error: err}
}

// HandleUpdateMovie godoc
//
//	@Summary		update a movie w/ ID
//	@Description	update a movie's details
//	@Tags			admin
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id				path		int							true	"movie ID to be updated"
//	@Param			title			formData	string						false	"movie title"
//	@Param			overview		formData	string						false	"movie description"
//	@Param			runtime			formData	int							false	"movie duration (minutes)"
//	@Param			backdrop_path	formData	file						false	"new backdrop file"
//	@Param			poster_path		formData	file						false	"new poster file"
//	@Success		200				{object}	models.UpdateMovieResponse	"movie updated successfully"
//	@Failure		400				{object}	models.UpdateMovieResponse	"invalid input or file"
//	@Failure		404				{object}	models.UpdateMovieResponse	"movie not found"
//	@Failure		500				{object}	models.UpdateMovieResponse	"server error"
//	@Security		BearerAuth
//	@Router			/admin/movies/{id} [patch]
func (m *MovieHandler) HandleMovieUpdate(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.PrintError("INVALID MOVIE ID", 12, err)
		ctx.JSON(http.StatusBadRequest, newUpdateMovieResponse(
			false, "", "invalid movie ID",
		))
		return
	}

	var newBody models.MovieBody
	if err := ctx.ShouldBind(&newBody); err != nil {
		utils.PrintError("UNABLE TO BIND MOVIE UPDATE BODY", 8, err)
		ctx.JSON(http.StatusInternalServerError, newUpdateMovieResponse(
			false, "", "server unable to bind body",
		))
		return
	}

	bDrop := newBody.NewBackdrop
	var bDropName string
	if bDrop != nil {
		ext := filepath.Ext(bDrop.Filename)
		filename := fmt.Sprintf("bdrop_%d_%d%s", idParam, time.Now().Unix(), ext)
		location := filepath.Join("public", "backdrop", filename)
		if err := ctx.SaveUploadedFile(bDrop, location); err != nil {
			utils.PrintError("INVALID MOVIE BACKDROP", 12, err)
			ctx.JSON(http.StatusBadRequest, newUpdateMovieResponse(
				false, "", "unable to upload movie backdrop",
			))
			return
		}
		bDropName = filename
		log.Println("====================", filename)
	}
	poster := newBody.NewPoster
	var posterName string
	if poster != nil {
		ext := filepath.Ext(poster.Filename)
		filename := fmt.Sprintf("poster_%d_%d%s", idParam, time.Now().Unix(), ext)
		location := filepath.Join("public", "poster", filename)
		if err := ctx.SaveUploadedFile(poster, location); err != nil {
			utils.PrintError("INVALID MOVIE POSTER", 12, err)
			ctx.JSON(http.StatusBadRequest, newUpdateMovieResponse(
				false, "", "unable to upload movie poster",
			))
			return
		}
		posterName = filename
		log.Println("====================", filename)
	}
	// log.Println("====================", newBody)

	if err := m.mr.UpdateMovie(newBody, bDropName, posterName, ctx.Request.Context(), idParam); err != nil {
		utils.PrintError("UNABLE TO UPDATE MOVIE TO DB", 8, err)
		ctx.JSON(http.StatusInternalServerError, newUpdateMovieResponse(
			false, "", "server unable to update movie",
		))
		return
	}

	ctx.JSON(http.StatusOK, newUpdateMovieResponse(
		true, fmt.Sprintf("movie w/ ID %d updated succesfully", idParam), "",
	))
}

// HandleGetMovieWithGenrePageSearch godoc
//
//	@Summary	get movie with filter by name and genre with pagination
//	@Tags		movies
//	@Accept		json
//	@Produce	json
//	@Param		q		query		string	false	"search title by q"	example(pulp)
//	@Param		page	query		int		false	"page number"		example(1)
//	@Param		genre	query		string	false	"genre"				example(action)
//	@Success	200		{object}	models.MovieResponse
//	@Router		/movies/ [get]
func (m *MovieHandler) HandleMovieWithGenrePageSearch(ctx *gin.Context) {
	q := ctx.Query("q")
	genreName := ctx.Query("genre")
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 12
	offset := (page - 1) * limit

	movies, err := m.mr.GetMovieWithGenrePageSearch(ctx.Request.Context(), q, genreName, limit, offset)
	if err != nil {
		utils.LogCtxError(
			ctx,
			"UNABLE TO GET MOVIE WITH FILTER",
			"Internal server error",
			err,
			http.StatusInternalServerError,
		)
		return
	}

	if len(movies) == 0 {
		ctx.JSON(http.StatusOK, newMoviesResponse(
			nil,
			true,
			"Nothing here but crickets",
		))
		return
	}
	ctx.JSON(http.StatusOK, newMoviesResponse(
		movies, true, "",
	))
}

func newMovieSchedulesResponse(res models.MovieSchedule, success bool, err string) models.MovieSchedulesResponse {
	return models.MovieSchedulesResponse{Result: res, Success: success, Error: err}
}

// HandleGetMovieSchedule godoc
//
//	@Summary		get movie schedules by ID
//	@Description	movie schedules (date, time, location, cinema) for a given movie ID
//	@Tags			schedule
//	@Produce		json
//	@Param			id	path		int								true	"Movie ID"
//	@Success		200	{object}	models.MovieSchedulesResponse	"Schedules retrieved successfully"
//	@Failure		500	{object}	models.MovieSchedulesResponse	"server error"
//	@Router			/movies/{id}/schedules [get]
func (m *MovieHandler) HandleGetMovieSchedule(ctx *gin.Context) {
	movieId, _ := strconv.Atoi(ctx.Param("id"))

	schedules, err := m.mr.GetMovieSchedules(ctx.Request.Context(), movieId)
	if err != nil {
		utils.PrintError(fmt.Sprintf("UNABLE TO GET MOVIE %d SCHEDULES", movieId), 12, err)
		ctx.JSON(http.StatusInternalServerError, newMovieSchedulesResponse(
			models.MovieSchedule{}, false, "server unable to get movie schedules",
		))
		return
	}
	if len(schedules.Schedule) == 0 {
		ctx.JSON(http.StatusNotFound, newMovieSchedulesResponse(
			models.MovieSchedule{}, false, fmt.Sprintf("no schedules found for movie ID %d", movieId),
		))
		return
	}

	ctx.JSON(http.StatusOK, newMovieSchedulesResponse(
		schedules, true, "",
	))
}

func newMovieScheduleFilterResponse(res []models.MovieScheduleFilter, success bool, err string) models.MovieScheduleFilterResponse {
	return models.MovieScheduleFilterResponse{Result: res, Success: success, Error: err}
}

// HandleGetMovieScheduleFilter godoc
//
//	@Summary		Get movie schedule filter
//	@Description	Retrieve available movie schedules based on movie ID, date, time, and location.
//	@Tags			schedule
//	@Param			id			path	int		true	"Movie ID"
//	@Param			date		query	string	true	"Filter by date (YYYY-MM-DD)"
//	@Param			time		query	int		true	"Filter by time (HHMM in 24h format, e.g. 2030)"
//	@Param			location	query	int		true	"Filter by location ID"
//	@Produce		json
//	@Success		200	{object}	models.MovieScheduleFilterResponse
//	@Failure		404	{object}	models.MovieScheduleFilterResponse	"No matching schedule found"
//	@Failure		500	{object}	models.MovieScheduleFilterResponse	"Internal server error"
//	@Router			/movies/{id}/schedule [get]
func (m *MovieHandler) HandleGetMovieScheduleFilter(ctx *gin.Context) {
	movieId, _ := strconv.Atoi(ctx.Param("id"))
	date := ctx.Query("date")
	time, _ := strconv.Atoi(ctx.Query("time"))
	loc, _ := strconv.Atoi(ctx.Query("location"))

	schedules, err := m.mr.GetMovieScheduleFilter(ctx.Request.Context(),
		movieId,
		time,
		loc,
		date,
	)
	if err != nil {
		utils.LogCtxError(
			ctx,
			"ERROR GET MOVIE FILTER SCHEDULE",
			"Internal server error",
			err,
			http.StatusInternalServerError,
		)
		return
	}

	if len(schedules) == 0 {
		ctx.JSON(http.StatusOK, newMovieScheduleFilterResponse(
			nil,
			true,
			"No matching movie schedule",
		))
		return
	}
	ctx.JSON(http.StatusOK, newMovieScheduleFilterResponse(
		schedules, true, "",
	))
}

func newCreateMovieResponse(res, err string, success bool) models.CreateMovieResponse {
	return models.CreateMovieResponse{Result: res, Success: success, Error: err}
}

// HandleCreateMovie godoc
//
//	@Summary		Create a new movie
//	@Description	Create a new movie with poster and backdrop upload, genres and casts as JSON arrays.
//	@Tags			admin
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id				formData	int							true	"Movie ID"
//	@Param			title			formData	string						true	"Movie title"
//	@Param			backdrop_path	formData	file						false	"Backdrop image file"
//	@Param			poster_path		formData	file						false	"Poster image file"
//	@Param			release_date	formData	string						true	"Release date in YYYY-MM-DD format"
//	@Param			runtime			formData	int							false	"Runtime in minutes"
//	@Param			overview		formData	string						false	"Movie overview"
//	@Param			director_name	formData	string						true	"Director's full name"
//	@Param			popularity		formData	number						true	"Popularity score (e.g. 78.5)"
//	@Param			genres			formData	string						false	"JSON array of genre IDs as string: [12,14,18]"
//	@Param			casts			formData	string						false	"JSON array of cast IDs as string: [1,2,3]"
//	@Success		200				{object}	models.CreateMovieResponse	"Successfully created movie"
//	@Failure		400				{object}	models.CreateMovieResponse	"Bad request, e.g. invalid input or file upload error"
//	@Failure		500				{object}	models.CreateMovieResponse	"Internal server error, e.g. binding failure"
//
//	@Security		BearerAuth
//
//	@Router			/admin/movies [post]
func (m *MovieHandler) HandleCreateMovie(ctx *gin.Context) {
	var body = models.CreateMovieBody{}

	if err := ctx.ShouldBind(&body); err != nil {
		utils.PrintError("UNABLE TO BIND CREATE MOVIE BODY", 12, err)
		ctx.JSON(http.StatusInternalServerError, newCreateMovieResponse(
			"", "server unable to bind request", false,
		))
		return
	}

	bDrop := body.BackdropPath
	var bDropName string
	if bDrop != nil {
		ext := filepath.Ext(bDrop.Filename)
		filename := fmt.Sprintf("bdrop_%d%s", time.Now().Unix(), ext)
		location := filepath.Join("public", "backdrop", filename)
		if err := ctx.SaveUploadedFile(bDrop, location); err != nil {
			utils.PrintError("INVALID MOVIE BACKDROP", 12, err)
			ctx.JSON(http.StatusBadRequest, newUpdateMovieResponse(
				false, "", "unable to upload movie backdrop",
			))
			return
		}
		bDropName = filename
		log.Println("====================", filename)
	}
	poster := body.PosterPath
	var posterName string
	if poster != nil {
		ext := filepath.Ext(poster.Filename)
		filename := fmt.Sprintf("poster_%d%s", time.Now().Unix(), ext)
		location := filepath.Join("public", "poster", filename)
		if err := ctx.SaveUploadedFile(poster, location); err != nil {
			utils.PrintError("INVALID MOVIE POSTER", 12, err)
			ctx.JSON(http.StatusBadRequest, newUpdateMovieResponse(
				false, "", "unable to upload movie poster",
			))
			return
		}
		posterName = filename
		log.Println("====================", filename)
	}
	// log.Println("====================", newBody)

	res, err := m.mr.CreateMovie(body, bDropName, posterName, ctx.Request.Context())
	if err != nil {
		utils.PrintError("ERROR CREATE MOVIE", 20, err)
		ctx.JSON(http.StatusInternalServerError, newCreateMovieResponse(
			"", "server error creating movie", true,
		))
		return
	}

	ctx.JSON(http.StatusOK, newCreateMovieResponse(
		fmt.Sprintf("success create movie w/ ID %d", res), "", true,
	))
}

func (m *MovieHandler) HandleGenres(ctx *gin.Context) {
	genres, err := m.mr.GetGenres(ctx)
	if err != nil {
		utils.LogCtxError(
			ctx,
			"SERVER UNABLE GET GENRES",
			"Internal server error",
			err,
			http.StatusInternalServerError,
		)
		return
	}

	ctx.JSON(http.StatusOK, models.NewFullfilledResponse(
		http.StatusOK,
		genres,
	))
}
