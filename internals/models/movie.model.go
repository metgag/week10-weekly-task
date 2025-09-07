package models

import (
	"time"
)

type Movie struct {
	ID           uint32    `db:"id" json:"id" example:"680"`
	Title        string    `db:"title" json:"title" example:"Pulp Fiction"`
	BackdropPath string    `db:"backdrop_path" json:"backdrop_path"`
	PosterPath   string    `db:"poster_path" json:"poster_path"`
	ReleaseDate  time.Time `db:"release_date" json:"release_date" example:"1994-09-10"`
	Runtime      uint16    `db:"runtime" json:"runtime" example:"154"`
	Overview     string    `db:"overview" json:"overview"`
	DirectorID   uint16    `db:"director_id" json:"director_id" example:"6"`
	Popularity   float32   `db:"popularity" json:"popularity" example:"17.246"`
}

type MovieResponse struct {
	Result  []Movie `json:"result"`
	Success bool    `json:"success"`
	Error   string  `json:"error"`
}

type MovieDetail struct {
	ID           uint32    `db:"id" json:"id" example:"680"`
	Title        string    `db:"title" json:"title" example:"Pulp Fiction"`
	BackdropPath string    `db:"backdrop_path" json:"backdrop_path"`
	PosterPath   string    `db:"poster_path" json:"poster_path"`
	ReleaseDate  time.Time `db:"release_date" json:"release_date" example:"1994-09-10"`
	Runtime      uint16    `db:"runtime" json:"runtime" example:"154"`
	Overview     string    `db:"overview" json:"overview"`
	DirectorName string    `db:"director_name" json:"director_name" example:"Quentin Tarantino"`
	Genres       []string  `db:"genres_name" json:"genres_name"`
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

type DeleteMovieResponse struct {
	Success bool
	Result  string `example:"movie w/ ID 1 deleted succesfully"`
	Error   string
}

type UpdateMovieResponse struct {
	Success bool
	Result  string `example:"movie w/ ID 1 updated succesfully"`
	Error   string
}
