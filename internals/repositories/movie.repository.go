package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/models"
)

type MovieRepository struct {
	dbpool *pgxpool.Pool
}

func NewMovieRepository(dbpool *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{dbpool: dbpool}
}

func (m *MovieRepository) UpcomingMoviesDat(ctx context.Context) ([]models.Movie, error) {
	sql := `
		SELECT 
			id, title, backdrop_path, poster_path, release_date, runtime, overview, director_id, popularity
		FROM 
			movies
		WHERE 
			release_date > current_date
	`
	rows, err := m.dbpool.Query(ctx, sql)
	if err != nil {
		return []models.Movie{}, nil
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.BackdropPath, &movie.PosterPath, &movie.ReleaseDate, &movie.Runtime, &movie.Overview, &movie.DirectorID, &movie.Popularity); err != nil {
			return []models.Movie{}, nil
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (m *MovieRepository) PopularMoviesDat(ctx context.Context) ([]models.Movie, error) {
	sql := `
		SELECT 
			id, title, backdrop_path, poster_path, release_date, runtime, overview, director_id, popularity
		FROM 
			movies
		WHERE 
			popularity > 40
		ORDER BY
			popularity DESC
	`
	rows, err := m.dbpool.Query(ctx, sql)
	if err != nil {
		return []models.Movie{}, nil
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.BackdropPath, &movie.PosterPath, &movie.ReleaseDate, &movie.Runtime, &movie.Overview, &movie.DirectorID, &movie.Popularity); err != nil {
			return []models.Movie{}, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (m *MovieRepository) MovieDetailsDat(ctx context.Context, id int) (models.MovieDetail, error) {
	sqlMovie := `
		SELECT
			m.id, m.title, m.poster_path, m.backdrop_path, m.release_date, m.runtime, m.overview, d.name "director"
		FROM
			movies AS m
		JOIN
			directors AS d ON m.director_id = d.id
		WHERE
			m.id = $1
	`
	sqlGenres := `
		SELECT
			g.genre_name
		FROM
			movies_genres AS m
		JOIN
			genres AS g
		ON
			m.genre_id = g.id
		WHERE
			m.movie_id = $1
	`
	rows, err := m.dbpool.Query(ctx, sqlGenres, id)
	if err != nil {
		return models.MovieDetail{}, err
	}
	defer rows.Close()

	var genres []string
	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			return models.MovieDetail{}, err
		}
		genres = append(genres, genre)
	}

	var movie models.MovieDetail
	if err := m.dbpool.QueryRow(ctx, sqlMovie, id).Scan(&movie.ID, &movie.Title, &movie.BackdropPath, &movie.PosterPath, &movie.ReleaseDate, &movie.Runtime, &movie.Overview, &movie.DirectorName); err != nil {
		return models.MovieDetail{}, err
	}
	movie.Genres = genres

	return movie, nil
}
