package models

type Response struct {
	Status  int
	Message string
	Data    map[string]interface{}
}
