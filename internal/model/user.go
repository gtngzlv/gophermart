package model

type User struct {
	ID       uint
	Login    string `json:"login"`
	Password string `json:"password"`
}
