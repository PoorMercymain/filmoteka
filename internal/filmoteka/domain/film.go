package domain

type Film struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Rating      *int   `json:"rating,omitempty"`
	Actors      []int  `json:"actorIDs"`
}
