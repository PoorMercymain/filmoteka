package domain

const (
	Male   = false
	Female = true
)

type Actor struct {
	Name     string `json:"name,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Birthday string `json:"birthday,omitempty"`
}

type OutputActor struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	Gender   string            `json:"gender"`
	Birthday string            `json:"birthday"`
	Films    []ActorOutputFilm `json:"films"`
}
