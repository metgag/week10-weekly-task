package models

import (
	"mime/multipart"
	"time"
)

type Movie struct {
	ID           uint32    `db:"id" json:"id" example:"680"`
	Title        string    `db:"title" json:"title" example:"Pulp Fiction" form:"title"`
	BackdropPath *string   `db:"backdrop_path" json:"backdrop_path"`
	PosterPath   *string   `db:"poster_path" json:"poster_path"`
	ReleaseDate  time.Time `db:"release_date" json:"release_date" example:"1994-09-10" form:"release_date"`
	Runtime      uint16    `db:"runtime" json:"runtime" example:"154" form:"runtime"`
	Overview     string    `db:"overview" json:"overview" form:"overview"`
	DirectorName string    `db:"director_name" json:"director_name" example:"Mick Jagger" form:"director_id"`
	// Popularity   *float32  `db:"popularity" json:"popularity" example:"17.246" form:"popularity"`
	Genres []Genre `db:"genres" json:"genres"`
	Casts  []Cast  `db:"casts" json:"cast"`
}

type MovieBody struct {
	Title       *string  `form:"title" db:"title"`
	Overview    *string  `form:"overview" db:"overview"`
	Runtime     *uint16  `form:"runtime" db:"runtime"`
	ReleaseDate *string  `form:"release_date" db:"release_date"`
	Popularity  *float32 `form:"popularity" db:"popularity"`
	// keep DirectorID if you want to update by ID directly
	DirectorID   *uint32               `form:"director_id" db:"director_id"`
	DirectorName *string               `form:"director_name"` // <-- NEW
	Casts        *string               `form:"casts"`         // CSV: "Tom Hanks, Leonardo DiCaprio"
	Genres       *string               `form:"genres"`        // CSV: "Drama, Thriller"
	NewBackdrop  *multipart.FileHeader `form:"backdrop_path"`
	NewPoster    *multipart.FileHeader `form:"poster_path"`
}

// type MovieBody struct {
// 	Title       *string               `form:"title" db:"title"`
// 	Overview    *string               `form:"overview" db:"overview"`
// 	Runtime     *uint16               `form:"runtime" db:"runtime"`
// 	ReleaseDate *string               `form:"release_date" db:"release_date"` // parse ke time.Time
// 	Popularity  *float32              `form:"popularity" db:"popularity"`
// 	DirectorID  *uint32               `form:"director_id" db:"director_id"`
// 	NewBackdrop *multipart.FileHeader `form:"backdrop_path"`
// 	NewPoster   *multipart.FileHeader `form:"poster_path"`
// }

// type MovieBody struct {
// 	Title       string                `form:"title"`
// 	Overview    string                `form:"overview"`
// 	Runtime     uint16                `form:"runtime"`
// 	ReleaseDate string                `form:"release_date"` // optionally parse to time.Time
// 	Popularity  float32               `form:"popularity"`
// 	DirectorID  uint32                `form:"director_id"` // gunakan ID, bukan nama
// 	NewBackdrop *multipart.FileHeader `form:"backdrop_path"`
// 	NewPoster   *multipart.FileHeader `form:"poster_path"`
// }

type MoviesResponse struct {
	Result  []MovieFilter `json:"result"`
	Success bool          `json:"success"`
	Error   string        `json:"error"`
}

type MovieResponse struct {
	Result  Movie  `json:"result"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// type MovieDetail struct {
// 	ID           uint32    `db:"id" json:"id" example:"680"`
// 	Title        string    `db:"title" json:"title" example:"Pulp Fiction"`
// 	BackdropPath string    `db:"backdrop_path" json:"backdrop_path"`
// 	PosterPath   string    `db:"poster_path" json:"poster_path"`
// 	ReleaseDate  time.Time `db:"release_date" json:"release_date" example:"1994-09-10"`
// 	Runtime      uint16    `db:"runtime" json:"runtime" example:"154"`
// 	Overview     string    `db:"overview" json:"overview"`
// 	DirectorName string    `db:"director_name" json:"director_name" example:"Quentin Tarantino"`
// 	Genres       *[]string `db:"genres_name" json:"genres_name"`
// }

// type MovieDetailResponse struct {
// 	Result  MovieDetail `json:"result"`
// 	Success bool        `json:"success"`
// 	Error   string      `json:"error"`
// }

type MovieFilter struct {
	ID          uint32    `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	PosterPath  *string   `db:"poster_path" json:"poster_path"`
	ReleaseDate time.Time `db:"release_date" json:"release_date"`
	Popularity  float32   `db:"popularity" json:"popularity"`
	Runtime     uint16    `db:"runtime" json:"runtime"`
	Genres      []Genre   `db:"genres" json:"genres"`
	Overview    string    `json:"overview"`
	Director    string    `json:"director"`
	Casts       string    `json:"casts"`
}

type MovieFilterResponse struct {
	Result  []MovieFilter `json:"result"`
	Success bool          `json:"success"`
	Error   string        `json:"error"`
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

type Genre struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type Cast struct {
	ID   uint16 `json:"id"`
	Name string `json:"name"`
}

type MovieSchedule struct {
	MovieID  uint32     `json:"movie_id"`
	Schedule []Schedule `json:"schedule"`
}

type Schedule struct {
	ScheduleID uint16    `json:"schedule_id"`
	Date       time.Time `json:"date"`
	Time       string    `json:"time"`
	Location   string    `json:"location"`
	Cinema     string    `json:"cinema"`
	CinemaImg  string    `json:"cinema_img"`
}

type MovieSchedulesResponse struct {
	Result  MovieSchedule `json:"result"`
	Success bool          `json:"success"`
	Error   string        `json:"error"`
}

type CreateMovieBody struct {
	DirectorName string                `form:"director_name" binding:"required"`
	Title        string                `form:"title" binding:"required"`
	BackdropPath *multipart.FileHeader `form:"backdrop_path"`
	PosterPath   *multipart.FileHeader `form:"poster_path"`
	ReleaseDate  string                `form:"release_date" binding:"required"`
	Runtime      uint16                `form:"runtime"`
	Overview     string                `form:"overview"`
	// Popularity   float32               `form:"popularity"`
	Genres       string `form:"genres"`
	Casts        string `form:"casts"`
	LocationID   []int  `form:"location"`
	ScheduleDate string `form:"schedule_date"`
	TimeID       []int  `form:"schedule_time"`
}

type CreateMovieResponse struct {
	Result  string `json:"result,omitempty"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type Director struct {
	ID   uint16 `db:"id"`
	Name string `db:"name"`
}

type MovieScheduleFilter struct {
	ScheduleID uint16 `json:"schedule_id"`
	CinemaID   uint16 `json:"cinema_id"`
	CinemaName string `json:"cinema_name"`
	CinemaImg  string `json:"cinema_img"`
}

type MovieScheduleFilterResponse struct {
	Result  []MovieScheduleFilter `json:"result"`
	Success bool                  `json:"success"`
	Error   string                `json:"error"`
}
