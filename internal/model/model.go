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

type GenerateReq struct {
	Values map[string]interface{} `json:"values"`
	Id     string                 `json:"-"`
}
