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

func newMovieResponse(result []models.Movie, success bool, error string) models.MovieResponse {
	return models.MovieResponse{Result: result, Success: success, Error: error}
}

func (m *MovieHandler) GetUpcomingMovies(ctx *gin.Context) {
	movies, err := m.mr.UpcomingMoviesDat(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newMovieResponse(
			[]models.Movie{}, false, "server unable to get movies",
		))
		return
	}

	if len(movies) == 0 {
		ctx.JSON(http.StatusNotFound, newMovieResponse(
			[]models.Movie{}, false, "movies data is empty",
		))
		return
	}

	ctx.JSON(http.StatusOK, newMovieResponse(
		movies, true, "",
	))
}

func (m *MovieHandler) GetPopularMovies(ctx *gin.Context) {
	movies, err := m.mr.PopularMoviesDat(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newMovieResponse(
			[]models.Movie{}, false, "server unable to get movies",
		))
		return
	}

	if len(movies) == 0 {
		ctx.JSON(http.StatusNotFound, newMovieResponse(
			[]models.Movie{}, false, "movies data is empty",
		))
		return
	}

	ctx.JSON(http.StatusOK, newMovieResponse(
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

func (m *MovieHandler) HandleGetAllMovie(ctx *gin.Context) {
	movies, err := m.mr.GetAllMoviesDat(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newMovieResponse(
			nil, false, "server unable to get list of movies",
		))
		return
	}
	if len(movies) == 0 {
		ctx.JSON(http.StatusInternalServerError, newMovieResponse(
			nil, false, "no movies on the server",
		))
		return
	}

	ctx.JSON(http.StatusOK, newMovieResponse(
		movies, true, "",
	))
}

func (m *MovieHandler) HandleDeleteMovie(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid movie id",
		})
		return
	}

	ctag, err := m.mr.DeleteMovie(ctx.Request.Context(), idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "server unable to delete movie",
		})
		return
	}

	if ctag.RowsAffected() == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("no movie w/ id: %d", idParam),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  ctag.String(),
	})
}

func (m *MovieHandler) HandleMovieUpdate(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid movie id",
		})
		return
	}

	var newBody models.Movie
	if err := ctx.ShouldBindJSON(&newBody); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "server unable to bind data",
		})
		return
	}

	ctag, err := m.mr.UpdateMovie(newBody, ctx.Request.Context(), idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "server unable to update movie data",
		})
		return
	}

	if ctag.RowsAffected() == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("unable to make update, no movie w/ ID %d", idParam),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  ctag.String(),
	})
}

func (m *MovieHandler) HandleMovieWithGenrePageSearch(ctx *gin.Context) {
	q := ctx.Query("q")
	genreId, err := strconv.Atoi(ctx.Query("genre_id"))
	if err != nil {
		genreId = 0
	}
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 4
	offset := (page - 1) * limit

	movies, err := m.mr.GetMovieWithGenrePageSearch(q, limit, offset, genreId, ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newMovieResponse(
			nil, false, "server unable to get movies data",
		))
		return
	}
	if len(movies) == 0 {
		ctx.JSON(http.StatusInternalServerError, newMovieResponse(
			nil, false, "no matching movie data on the server",
		))
		return
	}

	ctx.JSON(http.StatusOK, struct {
		Result  []models.MovieGenre `json:"result"`
		Success bool                `json:"success"`
	}{
		movies,
		true,
	})
}
