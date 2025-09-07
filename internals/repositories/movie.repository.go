package repositories

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/jackc/pgx/v5/pgconn"
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

func (m *MovieRepository) GetAllMoviesDat(ctx context.Context) ([]models.Movie, error) {
	sql := `
		SELECT
			id, title, backdrop_path, poster_path, release_date, runtime, overview, director_id, popularity
		FROM
			movies
	`
	rows, err := m.dbpool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.BackdropPath,
			&movie.PosterPath,
			&movie.ReleaseDate,
			&movie.Runtime,
			&movie.Overview,
			&movie.DirectorID,
			&movie.Popularity,
		); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (m *MovieRepository) DeleteMovie(ctx context.Context, id int) (pgconn.CommandTag, error) {
	sqlMovie := `
		DELETE FROM movies 
		WHERE id = $1
	`
	sqlGetGenre := `
		SELECT movie_id
		FROM movies_genres
		WHERE movie_id = $1
	`
	ctagGetGenre, _ := m.dbpool.Exec(ctx, sqlGetGenre, id)
	log.Println(ctagGetGenre.RowsAffected())
	if ctagGetGenre.RowsAffected() == 0 {
		return m.dbpool.Exec(ctx, sqlMovie, id)
	}

	sqlDeleteGenre := `
		DELETE FROM movies_genres
		WHERE movie_id = $1
	`
	_, err := m.dbpool.Exec(ctx, sqlDeleteGenre, id)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	return m.dbpool.Exec(ctx, sqlMovie, id)
}

func (m *MovieRepository) UpdateMovie(newBody models.Movie, ctx context.Context, id int) (pgconn.CommandTag, error) {
	rt := reflect.TypeOf(newBody)
	rv := reflect.ValueOf(newBody)

	var args []any
	var argIndex int = 1

	sql := "UPDATE movies SET "
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		if value.IsZero() {
			continue
		} else {
			args = append(args, value.Interface())
		}

		sql += fmt.Sprintf("%s = $%d", field.Tag.Get("db"), argIndex)
		sql += ", "

		argIndex++
	}

	sql += fmt.Sprintf(" updated_at = current_timestamp WHERE id = $%d", argIndex)
	args = append(args, id)

	return m.dbpool.Exec(ctx, sql, args...)
}

func (m *MovieRepository) GetMovieWithGenrePageSearch(q string, limit, offset, genreId int, ctx context.Context) ([]models.MovieGenre, error) {
	sql := `
		SELECT 
			m.id, m.title, mg.genre_id
		FROM
			movies AS m
		JOIN
			movies_genres AS mg ON m.id = mg.movie_id
		WHERE
			mg.genre_id = $1
		AND
			m.title ILIKE $2
		LIMIT $3 OFFSET $4
	`
	search := "%" + q + "%"
	rows, err := m.dbpool.Query(ctx, sql, genreId, search, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieGenre
	for rows.Next() {
		var movie models.MovieGenre
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.GenreID); err != nil {
			return nil, err
		}
		log.Println("fofofofoffofofoofofofofoofofofofoff")
		movies = append(movies, movie)
	}

	return movies, nil
}
