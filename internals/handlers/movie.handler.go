package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/repositories"
)

type MovieHandler struct {
	mr *repositories.MovieRepository
}

func NewMovieHandler(mr *repositories.MovieRepository) *MovieHandler {
	return &MovieHandler{mr: mr}
}

func NewMovieResponse(result []models.Movie, success bool, error string) models.MovieResponse {
	return models.MovieResponse{Result: result, Success: success, Error: error}
}

func (m *MovieHandler) GetUpcomingMovies(ctx *gin.Context) {
	movies, err := m.mr.UpcomingMoviesDat(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewMovieResponse(
			[]models.Movie{}, false, "server unable to get movies",
		))
		return
	}

	if len(movies) == 0 {
		ctx.JSON(http.StatusNotFound, NewMovieResponse(
			[]models.Movie{}, false, "movies data is empty",
		))
		return
	}

	ctx.JSON(http.StatusOK, NewMovieResponse(
		movies, true, "",
	))
}

func (m *MovieHandler) GetPopularMovies(ctx *gin.Context) {
	movies, err := m.mr.PopularMoviesDat(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewMovieResponse(
			[]models.Movie{}, false, "server unable to get movies",
		))
		return
	}

	if len(movies) == 0 {
		ctx.JSON(http.StatusNotFound, NewMovieResponse(
			[]models.Movie{}, false, "movies data is empty",
		))
		return
	}

	ctx.JSON(http.StatusOK, NewMovieResponse(
		movies, true, "",
	))
}

func NewMovieDetailResponse(res models.MovieDetail, success bool, err string) models.MovieDetailResponse {
	return models.MovieDetailResponse{Result: res, Success: success, Error: err}
}

func (m *MovieHandler) GetMovieDetail(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, NewMovieDetailResponse(
			models.MovieDetail{}, false, "invalid movie id input",
		))
		return
	}

	detail, err := m.mr.MovieDetailsDat(ctx.Request.Context(), idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewMovieDetailResponse(
			models.MovieDetail{}, false, fmt.Sprintf("server unable to reach movie w/ ID %d", idParam),
		))
		return
	}

	ctx.JSON(http.StatusOK, NewMovieDetailResponse(
		detail, true, "",
	))
}
