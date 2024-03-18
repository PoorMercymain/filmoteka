package domain

type Film struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	ReleaseDate string   `json:"releaseDate,omitempty"`
	Rating      *float32 `json:"rating,omitempty"`
	Actors      []int    `json:"actorIDs"`
}

type OutputFilm struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ReleaseDate string  `json:"releaseDate"`
	Rating      float32 `json:"rating"`
	ActorIDs    []int   `json:"actorIDs"`
}

type ActorOutputFilm struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ReleaseDate string  `json:"releaseDate"`
	Rating      float32 `json:"rating"`
}