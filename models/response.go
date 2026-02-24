package models

type APIResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
}
