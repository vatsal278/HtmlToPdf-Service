package model

type PingRequest struct {
	Data string `json:"data" validate:"required"`
}

type Student struct {
	Name  string
	Marks int
	Id    string
}
type Class []Student

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    interface{}
}
