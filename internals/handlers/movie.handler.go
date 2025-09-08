package handlers

import (
	"fmt"
	"log"
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

// HandleGetUpcomingMovies godoc
//
//	@Summary	get upcoming movies handler func
//	@Tags		movies
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.MovieResponse
//	@Router		/movies/upcoming [get]
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

// HandleGetPopularMovies godoc
//
//	@Summary	get popular movies handler func
//	@Tags		movies
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.MovieResponse
//	@Router		/movies/popular [get]
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

// HandleGetMovieDetail godoc
//
//	@Summary	get movie detail handler func
//	@Tags		movies
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"movie ID"
//	@Success	200	{object}	models.MovieDetailResponse
//	@Router		/movies/{id} [get]
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

// HandleGetAllMovies godoc
//
//	@Summary	get all movies (admin)
//	@Tags		movies
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Success	200				{object}	models.MovieResponse
//	@Router		/admin/movies/ [get]
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

func newDeleteMovieResponse(success bool, res, err string) models.DeleteMovieResponse {
	return models.DeleteMovieResponse{Success: success, Result: res, Error: err}
}

// HandleDeleteMovie godoc
//
//	@Summary	delete a movie w/ ID
//	@Tags		movies
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer token"
//	@Param		id				path		int		true	"movie ID"
//	@Success	200				{object}	models.DeleteMovieResponse
//	@Router		/admin/movies/{id} [delete]
func (m *MovieHandler) HandleDeleteMovie(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newDeleteMovieResponse(
			false, "", "invalid movie ID",
		))
		return
	}

	ctag, err := m.mr.DeleteMovie(ctx.Request.Context(), idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newDeleteMovieResponse(
			false, "", "server unable to delete movie",
		))
		return
	}

	if ctag.RowsAffected() == 0 {
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
//	@Summary	update a movie w/ ID
//	@Tags		movies
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string			true	"Bearer token"
//	@Param		id				path		int				true	"movie ID"
//	@Param		request			body		models.Movie	true	"movie body json content"
//	@Success	200				{object}	models.UpdateMovieResponse
//	@Router		/admin/movies/{id} [patch]
func (m *MovieHandler) HandleMovieUpdate(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Fatal(err.Error())
		ctx.JSON(http.StatusBadRequest, newUpdateMovieResponse(
			false, "", "invalid movie ID",
		))
		return
	}

	var newBody models.Movie
	if err := ctx.ShouldBindJSON(&newBody); err != nil {
		ctx.JSON(http.StatusInternalServerError, newUpdateMovieResponse(
			false, "", "server unable to bind body",
		))
		return
	}

	ctag, err := m.mr.UpdateMovie(newBody, ctx.Request.Context(), idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newUpdateMovieResponse(
			false, "", "server unable to update movie",
		))
		return
	}

	if ctag.RowsAffected() == 0 {
		ctx.JSON(http.StatusInternalServerError, newUpdateMovieResponse(
			false, "", fmt.Sprintf("no movie w/ ID %d", idParam),
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
//	@Param		example	query		string	false	"search title by q"	example(gump)
//	@Param		int		query		int		false	"page number"		example(1)
//	@Param		default	query		int		false	"genre_id default"	default(0)
//	@Param		enumint	query		int		false	"genre_id enums"	Enums(27, 53, 12, 14, 28, 80, 18, 878, 35, 10749)
//	@Success	200		{object}	models.MovieResponse
//	@Router		/movies/ [get]
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
