package domain

type AuthorizationData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}