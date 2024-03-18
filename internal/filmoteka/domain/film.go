package domain

type Film struct {
	Title       string   `json:"title,omitempty" example:"film 2"`
	Description string   `json:"description,omitempty" example:"some kind of film"`
	ReleaseDate string   `json:"releaseDate,omitempty" example:"2007-09-20"`
	Rating      *float32 `json:"rating,omitempty" example:"8.6"`
	Actors      []int    `json:"actorIDs" example:"1,2,3"`
}

type OutputFilm struct {
	ID          int               `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	ReleaseDate string            `json:"releaseDate"`
	Rating      float32           `json:"rating"`
	Actors      []FilmOutputActor `json:"actors"`
}

type ActorOutputFilm struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ReleaseDate string  `json:"releaseDate"`
	Rating      float32 `json:"rating"`
}
