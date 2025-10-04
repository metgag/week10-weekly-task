package repositories

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/models"
	"github.com/metgag/koda-weekly10/internals/utils"
	"github.com/redis/go-redis/v9"
)

type MovieRepository struct {
	dbpool *pgxpool.Pool
	rdb    *redis.Client
}

func NewMovieRepository(dbpool *pgxpool.Pool, rdb *redis.Client) *MovieRepository {
	return &MovieRepository{dbpool: dbpool, rdb: rdb}
}

func (m *MovieRepository) GetMovieSchedules(ctx context.Context, movieId int) (models.MovieSchedule, error) {
	sql := `
		SELECT
			cs.id, cs.show_date, t.show_time, l.show_location, c.cinema_name, c.cinema_img
		FROM
			schedule cs
		JOIN
			movies m ON m.id = cs.movie_id
		JOIN
			jam_tayang t ON t.id = cs.time_id
		JOIN
			lokasi_tayang l ON l.id = cs.location_id
		JOIN
			cinema_tayang c ON c.id = cs.cinema_id
		WHERE
			m.id = $1
		ORDER BY
			cs.id ASC
	`
	rows, err := m.dbpool.Query(ctx, sql, movieId)
	if err != nil {
		return models.MovieSchedule{}, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(
			&s.ScheduleID,
			&s.Date,
			&s.Time,
			&s.Location,
			&s.Cinema,
			&s.CinemaImg,
		); err != nil {
			return models.MovieSchedule{}, err
		}
		schedules = append(schedules, s)
	}

	movieSchedule := models.MovieSchedule{
		MovieID:  uint32(movieId),
		Schedule: schedules,
	}
	return movieSchedule, nil
}

func (m *MovieRepository) GetMovieScheduleFilter(ctx context.Context, movieId, timeId, locationId int, date string) ([]models.MovieScheduleFilter, error) {
	sql := `
		SELECT s.id, ct.id, ct.cinema_name, ct.cinema_img
		FROM schedule s
		JOIN cinema_tayang ct ON ct.id = s.cinema_id
		WHERE s.movie_id = $1
		AND s.show_date = $2
		AND s.time_id = $3
		AND s.location_id = $4
	`
	rows, err := m.dbpool.Query(ctx, sql, movieId, date, timeId, locationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.MovieScheduleFilter
	for rows.Next() {
		var schedule models.MovieScheduleFilter
		if err := rows.Scan(
			&schedule.ScheduleID,
			&schedule.CinemaID,
			&schedule.CinemaName,
			&schedule.CinemaImg,
		); err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (m *MovieRepository) UpdateMovie(
	newBody models.MovieBody,
	backdropPath string,
	posterPath string,
	ctx context.Context,
	id int,
) (pgconn.CommandTag, error) {
	tx, err := m.dbpool.Begin(ctx)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	rt := reflect.TypeOf(newBody)
	rv := reflect.ValueOf(newBody)

	var setClauses []string
	var args []any
	argIndex := 1

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		dbTag := field.Tag.Get("db")
		if dbTag == "" || value.IsZero() {
			continue
		}

		setClause := fmt.Sprintf("%s = $%d", dbTag, argIndex)
		setClauses = append(setClauses, setClause)
		args = append(args, value.Interface())
		argIndex++
	}

	if posterPath != "" {
		setClauses = append(setClauses, fmt.Sprintf("poster_path = $%d", argIndex))
		args = append(args, posterPath)
		argIndex++
	}

	if backdropPath != "" {
		setClauses = append(setClauses, fmt.Sprintf("backdrop_path = $%d", argIndex))
		args = append(args, backdropPath)
		argIndex++
	}

	setClauses = append(setClauses, "updated_at = current_timestamp")

	// Final query
	sql := fmt.Sprintf("UPDATE movies SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIndex)
	args = append(args, id)

	return tx.Exec(ctx, sql, args...)
}

func (m *MovieRepository) GetMovieWithGenrePageSearch(ctx context.Context, q, genreName string, limit, offset int) ([]models.MovieFilter, error) {
	baseSQL := `
		SELECT DISTINCT
			m.id, m.title, m.poster_path, m.release_date, m.runtime
		FROM
			movies m
		JOIN
			movies_genres mg ON mg.movie_id = m.id
		WHERE
			m.title ILIKE $1
	`

	args := []any{"%" + q + "%"} // $1
	if genreName != "" {
		genreID, err := m.castGenre(genreName)
		if err != nil {
			return nil, err
		}
		if genreID != 0 {
			baseSQL += " AND mg.genre_id = $2"
			args = append(args, genreID) // $2
		}
	}

	paramIdx := len(args) + 1
	baseSQL += fmt.Sprintf(" AND deleted_at IS NULL ORDER BY m.id ASC LIMIT $%d OFFSET $%d", paramIdx, paramIdx+1)
	args = append(args, limit, offset)

	rows, err := m.dbpool.Query(ctx, baseSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieFilter
	for rows.Next() {
		var movie models.MovieFilter
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.PosterPath,
			&movie.ReleaseDate,
			&movie.Runtime,
		); err != nil {
			return nil, err
		}

		genres, err := m.fetchGenres(ctx, int(movie.ID))
		if err != nil {
			return nil, err
		}
		movie.Genres = genres

		movies = append(movies, movie)
	}

	return movies, nil
}

func (m *MovieRepository) SoftDeleteMovie(ctx context.Context, movieId int) (pgconn.CommandTag, error) {
	sql := `
		UPDATE 
			movies
		SET
			is_deleted = true,
			deleted_at = current_timestamp
		WHERE
			id = $1
	`

	return m.dbpool.Exec(ctx, sql, movieId)
}

func (m *MovieRepository) GetPopularMovies(ctx context.Context) ([]models.MovieFilter, error) {
	redisKey := "archie:movies_populars"
	var populars []models.MovieFilter

	isExist, err := utils.CacheGet(m.rdb, ctx, redisKey, &populars)
	if err != nil {
		utils.PrintError("redis> REDIS ERROR", 20, err)
	}
	if isExist {
		return populars, nil
	}

	sql := `
		SELECT m.id, m.title, m.poster_path, m.release_date
		FROM orders o
		JOIN schedule s ON s.id = o.schedule_id
		JOIN movies m ON m.id = s.movie_id
		GROUP BY m.id, m.title, m.poster_path, m.release_date
		ORDER BY COUNT(m.id) DESC
	`

	rows, err := m.dbpool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.MovieFilter
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.PosterPath,
			&movie.ReleaseDate,
		); err != nil {
			return nil, err
		}

		genres, err := m.fetchGenres(ctx, int(movie.ID))
		if err != nil {
			return nil, err
		}

		movie.Genres = append(movie.Genres, genres...)
		populars = append(populars, movie)
	}

	expiration := 24 * time.Hour
	if err := utils.CacheSet(m.rdb, ctx, redisKey, populars, expiration); err != nil {
		utils.PrintError(
			fmt.Sprintf("redis> UNABLE TO SET %s", redisKey), 20, err,
		)
	}

	return populars, nil
}

func (m *MovieRepository) GetUpcomingMovies(ctx context.Context) ([]models.MovieFilter, error) {
	redisKey := "archie:movies_upcomings"
	var cached []models.MovieFilter

	isExist, err := utils.CacheGet(m.rdb, ctx, redisKey, &cached)
	if err != nil {
		utils.PrintError("redis> REDIS ERROR", 20, err)
	}
	if isExist {
		return cached, nil
	}

	// var upcomings []models.MovieFilter
	upcomings, err := m.GetAllMovies(ctx, `
		WHERE release_date + INTERVAL '1 months' > CURRENT_DATE
	`)
	if err != nil {
		return nil, err
	}

	expiration := 24 * time.Hour
	if err := utils.CacheSet(m.rdb, ctx, redisKey, upcomings, expiration); err != nil {
		utils.PrintError(
			fmt.Sprintf("redis> UNABLE TO SET %s", redisKey), 20, err,
		)
	}

	return upcomings, nil
}

func (m *MovieRepository) GetAllMovies(ctx context.Context, opts string) ([]models.MovieFilter, error) {
	sql := `
		SELECT 
			id, title, poster_path, release_date, runtime
		FROM
			movies
	`
	if opts != "" {
		sql += opts
		sql += "AND"
	} else {
		sql += "WHERE"
	}
	sql += `
			deleted_at IS NULL
		ORDER BY
			id ASC
	`
	// utils.PrintError(fmt.Sprintf("MOVIES QUERY\n%s", sql), 20, nil)

	rows, err := m.dbpool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieFilter
	for rows.Next() {
		var movie models.MovieFilter
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.PosterPath,
			&movie.ReleaseDate,
			&movie.Runtime,
		); err != nil {
			return nil, err
		}

		genres, err := m.fetchGenres(ctx, int(movie.ID))
		if err != nil {
			return nil, err
		}
		movie.Genres = append(movie.Genres, genres...)

		movies = append(movies, movie)
	}

	utils.PrintError(fmt.Sprintf("SQL QUERY \n%s", sql), 20, nil)

	return movies, nil
}

func (m *MovieRepository) GetMovieDetail(ctx context.Context, movieId int) (models.Movie, error) {
	sql := `
		SELECT
			m.id, m.title, m.backdrop_path, m.poster_path, m.release_date, m.runtime, m.overview, d.name
		FROM 
			movies m
		JOIN 
			directors d ON d.id = m.director_id
		WHERE
			m.id = $1
	`

	var movie models.Movie
	if err := m.dbpool.QueryRow(ctx, sql, movieId).Scan(
		&movie.ID,
		&movie.Title,
		&movie.BackdropPath,
		&movie.PosterPath,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Overview,
		&movie.DirectorName,
	); err != nil {
		return models.Movie{}, err
	}

	genres, err := m.fetchGenres(ctx, movieId)
	if err != nil {
		return models.Movie{}, err
	}
	casts, err := m.fetchCasts(ctx, movieId)
	if err != nil {
		return models.Movie{}, err
	}
	movie.Genres = append(movie.Genres, genres...)
	movie.Casts = append(movie.Casts, casts...)

	return movie, nil
}

func (m *MovieRepository) fetchCasts(ctx context.Context, movieId int) ([]models.Cast, error) {
	sql := `
		SELECT a.id, a.name
		FROM movies_casts ma
		JOIN casts a ON a.id = ma.cast_id
		JOIN movies m ON m.id = ma.movie_id
		WHERE m.id = $1;
	`
	rows, err := m.dbpool.Query(ctx, sql, movieId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var casts []models.Cast
	for rows.Next() {
		var cast models.Cast
		if err := rows.Scan(&cast.ID, &cast.Name); err != nil {
			return nil, err
		}

		casts = append(casts, cast)
	}

	return casts, nil
}

func (m *MovieRepository) fetchGenres(ctx context.Context, movieId int) ([]models.Genre, error) {
	sql := `
		SELECT g.id, g.genre_name
		FROM movies_genres mg
		JOIN genres g ON g.id = mg.genre_id
		JOIN movies m ON m.id = mg.movie_id
		WHERE m.id = $1;
	`
	rows, err := m.dbpool.Query(ctx, sql, movieId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}

	return genres, nil
}

func (m *MovieRepository) castGenre(strGenre string) (int, error) {
	switch strings.ToLower(strGenre) {
	case "action":
		return 28, nil
	case "adventure":
		return 12, nil
	case "animation":
		return 16, nil
	case "comedy":
		return 35, nil
	case "crime":
		return 80, nil
	case "documentary":
		return 99, nil
	case "drama":
		return 18, nil
	case "family":
		return 10751, nil
	case "fantasy":
		return 14, nil
	case "history":
		return 36, nil
	case "horror":
		return 27, nil
	case "music":
		return 10402, nil
	case "mystery":
		return 9648, nil
	case "romance":
		return 10749, nil
	case "sci-fi", "science-fiction", "science fiction", "scifi":
		return 878, nil
	case "tv-movie", "tvmovie", "tv":
		return 10770, nil
	case "thriller":
		return 53, nil
	case "war":
		return 10752, nil
	case "western":
		return 37, nil
	default:
		return 0, fmt.Errorf("invalid genre: %q", strGenre)
	}
}

func (m *MovieRepository) CreateMovie(
	body models.CreateMovieBody,
	backdropPath string,
	posterPath string,
	ctx context.Context,
) (uint32, error) {
	tx, err := m.dbpool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	rt := reflect.TypeOf(body)
	rv := reflect.ValueOf(body)

	var columns []string
	var placeholders []string
	var args []any
	argIndex := 1

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		formTag := field.Tag.Get("form")
		if formTag == "" || value.IsZero() {
			continue
		}

		if formTag == "director_name" {
			directorId, err := m.insertDirectors(tx, ctx, value.String())
			if err != nil {
				return 0, err
			}
			columns = append(columns, "director_id")
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, directorId)
			argIndex++
			continue
		}

		if field.Type.String() == "*multipart.FileHeader" || formTag == "genres" || formTag == "casts" {
			continue
		}

		if formTag == "release_date" {
			dateStr := value.String()
			dateVal, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return 0, fmt.Errorf("invalid date format for release_date")
			}
			columns = append(columns, "release_date")
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, dateVal)
			argIndex++
			continue
		}

		if formTag == "location" ||
			formTag == "schedule_date" ||
			formTag == "schedule_time" {
			continue
		}

		columns = append(columns, formTag)
		placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
		args = append(args, value.Interface())
		argIndex++
	}

	if posterPath != "" {
		columns = append(columns, "poster_path")
		placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
		args = append(args, posterPath)
		argIndex++
	}

	if backdropPath != "" {
		columns = append(columns, "backdrop_path")
		placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
		args = append(args, backdropPath)
		argIndex++
	}

	sql := fmt.Sprintf(
		"INSERT INTO movies (%s) VALUES (%s) RETURNING id",
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	// log.Println("Executing SQL:", sql)

	var newMovieID uint32
	if err := tx.QueryRow(ctx, sql, args...).Scan(&newMovieID); err != nil {
		return 0, err
	}

	// insert genres
	if _, err := m.insertMovieGenres(tx, ctx, newMovieID, body.Genres); err != nil {
		return 0, err
	}
	// insert casts
	// log.Printf("Insert casts for movieID %d: %s", newMovieID, body.Casts)
	log.Println(body.Casts)
	if err := m.insertMovieCasts(tx, ctx, newMovieID, body.Casts); err != nil {
		return 0, err
	}
	if err := m.createMovieSchedule(tx, ctx, newMovieID, body.ScheduleDate, body.LocationID, body.TimeID); err != nil {
		return 0, err
	}
	// log.Println("Insert casts succeeded")

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	redisKey := "archie:movies_upcomings"
	res, err := m.rdb.Del(ctx, redisKey).Result()
	if err != nil {
		log.Println(err)
	}
	log.Printf("Number of keys deleted: %d", res)

	return newMovieID, nil
}

func (m *MovieRepository) createMovieSchedule(
	tx pgx.Tx,
	ctx context.Context,
	movieId uint32,
	scheduleDate string,
	locationId,
	timeId []int,
) error {
	sql := `
		INSERT INTO 
			schedule(movie_id, show_date, time_id, location_id, cinema_id)
		VALUES
			($1, $2, $3, $4, $5)
		ON CONFLICT 
			(movie_id, show_date, time_id, location_id, cinema_id)
		DO NOTHING
	`
	for _, tId := range timeId {
		for _, lId := range locationId {
			for _, cId := range m.sliceCinemaID(lId) {
				if _, err := tx.Exec(
					ctx,
					sql,
					movieId,
					scheduleDate,
					tId,
					lId,
					cId,
				); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *MovieRepository) sliceCinemaID(locationId int) []int {
	switch locationId {
	case 1:
		return []int{1, 2, 3, 4}
	case 2:
		return []int{5, 6, 7, 8}
	case 3:
		return []int{9, 10, 11, 12}
	default:
		return []int{}
	}
}

// func (m* MovieRepository) createMovieLocationSchedule(
// 	tx pgx.Tx,
// 	ctx context.Context,
// ) {
// 	sql := `

// 	`
// }

func (m *MovieRepository) insertMovieCasts(tx pgx.Tx, ctx context.Context, movieID uint32, castCSV string) error {
	castStrs := strings.SplitSeq(castCSV, ",")

	for str := range castStrs {
		sqlC := `
			WITH inse AS (
			INSERT INTO casts(name)
			VALUES ($1)
			ON CONFLICT (name) DO NOTHING
			RETURNING id
			)
			SELECT id FROM inse
			UNION
			SELECT id FROM casts WHERE name = $1
			LIMIT 1;
		`
		var castID uint16
		if err := tx.QueryRow(ctx, sqlC, str).Scan(&castID); err != nil {
			return err
		}

		sqlMC := `
			INSERT INTO
				movies_casts(movie_id, cast_id)
			VALUES
				($1, $2)
		`

		_, err := tx.Exec(ctx, sqlMC, movieID, castID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MovieRepository) insertMovieGenres(tx pgx.Tx, ctx context.Context, movieID uint32, genreCSV string) (pgconn.CommandTag, error) {
	if strings.TrimSpace(genreCSV) == "" {
		return pgconn.CommandTag{}, nil
	}

	// Split and clean genre strings
	genreStrs := strings.Split(genreCSV, ",")
	var genreIDs []int

	for _, str := range genreStrs {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}

		id, err := m.castGenre(str)
		if err != nil {
			// Skip unknown/invalid genre (optional: log the error)
			continue
		}
		genreIDs = append(genreIDs, id)
	}

	if len(genreIDs) == 0 {
		return pgconn.CommandTag{}, nil
	}

	// Prepare placeholders and args
	var (
		placeholders []string
		args         []any
		argIndex     = 1
	)

	for _, genreID := range genreIDs {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", argIndex, argIndex+1))
		args = append(args, movieID, genreID)
		argIndex += 2
	}

	query := fmt.Sprintf(
		"INSERT INTO movies_genres (movie_id, genre_id) VALUES %s",
		strings.Join(placeholders, ", "),
	)

	ctag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to insert movie genres: %w", err)
	}

	return ctag, nil
}

func (m *MovieRepository) insertDirectors(tx pgx.Tx, ctx context.Context, directorName string) (uint16, error) {
	sql := `
		WITH ins AS (
		INSERT INTO directors(name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
		RETURNING id
		)
		SELECT id FROM ins
		UNION
		SELECT id FROM directors WHERE name = $1
		LIMIT 1;
	`
	var id uint16
	if err := tx.QueryRow(ctx, sql, directorName).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (m *MovieRepository) GetGenres(ctx context.Context) ([]models.Genre, error) {
	sql := `
		SELECT id, genre_name
		FROM genres
	`

	rows, err := m.dbpool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre

	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(
			&genre.ID,
			&genre.Name,
		); err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}

	return genres, nil
}
