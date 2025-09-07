package models

import "github.com/jackc/pgx/v5/pgtype"

type Movie struct {
	ID           uint32      `db:"id" json:"id"`
	Title        string      `db:"title" json:"title"`
	BackdropPath string      `db:"backdrop_path" json:"backdrop_path"`
	PosterPath   string      `db:"poster_path" json:"poster_path"`
	ReleaseDate  pgtype.Date `db:"release_date" json:"release_date"`
	Runtime      uint16      `db:"runtime" json:"runtime"`
	Overview     string      `db:"overview" json:"overview"`
	DirectorID   uint16      `db:"director_id" json:"director_id"`
	Popularity   float32     `db:"popularity" json:"popularity"`
}

type MovieResponse struct {
	Result  []Movie `json:"result"`
	Success bool    `json:"success"`
	Error   string  `json:"error"`
}

type MovieDetail struct {
	ID           uint32      `db:"id" json:"id"`
	Title        string      `db:"title" json:"title"`
	BackdropPath string      `db:"backdrop_path" json:"backdrop_path"`
	PosterPath   string      `db:"poster_path" json:"poster_path"`
	ReleaseDate  pgtype.Date `db:"release_date" json:"release_date"`
	Runtime      uint16      `db:"runtime" json:"runtime"`
	Overview     string      `db:"overview" json:"overview"`
	DirectorName string      `db:"director_name" json:"director_name"`
	Genres       []string    `db:"genres_name" json:"genres_name"`
}

type MovieDetailResponse struct {
	Result  MovieDetail `json:"result"`
	Success bool        `json:"success"`
	Error   string      `json:"error"`
}

type MovieGenre struct {
	ID      uint32 `db:"movie_id" json:"movie_id"`
	Title   string `db:"title" json:"title"`
	GenreID uint16 `db:"genre_id" json:"genre_id"`
}
