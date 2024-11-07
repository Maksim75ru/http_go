package main

import (
	"net/http"
)

func main() {
	memoryStorage := NewMemoryStorage()
	handler := NewHandler(memoryStorage)

	http.HandleFunc("POST employee/create", handler.CreateEmployee)

	http.ListenAndServe(":8080", nil)
}
