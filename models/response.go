package models

type Reply struct {
	Error   string
	Message string
	Success bool
	Data    interface{}
}
