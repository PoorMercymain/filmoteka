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
