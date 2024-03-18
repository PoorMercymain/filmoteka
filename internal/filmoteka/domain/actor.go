package domain

const (
	Male   = false
	Female = true
)

type Actor struct {
	Name     string `json:"name,omitempty" example:"Vasily Abcd"`
	Gender   string `json:"gender,omitempty" example:"male"`
	Birthday string `json:"birthday,omitempty" example:"2001-10-25"`
}

type OutputActor struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	Gender   string            `json:"gender"`
	Birthday string            `json:"birthday"`
	Films    []ActorOutputFilm `json:"films"`
}

type FilmOutputActor struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}
