package main

import (
	"net/http"
)

func main() {
	memoryStorage := NewMemoryStorage()
	handler := NewHandler(memoryStorage)

	http.HandleFunc("POST /employee/create", handler.CreateEmployee)
	http.HandleFunc("GET /employees", handler.GetEmployees)
	http.HandleFunc("GET /employee/{id}", handler.GetEmployee)
	http.HandleFunc("PUT /employee/{id}", handler.UpdateEmployee)
	http.HandleFunc("DELETE /employee/{id}", handler.DeleteEmployee)

	http.ListenAndServe(":8080", nil)
}
