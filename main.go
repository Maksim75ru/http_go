package main

import (
	"http_go/storages"
	"net/http"
)

func main() {
	// Ниже рабочая часть с MongoDb
	mongoDbStorage := storages.NewMongoDbStorage("RogaAndKopita", "employees")
	handler := NewHandler(mongoDbStorage)

	http.HandleFunc("POST /employees/create", handler.CreateEmployee)
	http.HandleFunc("GET /employees", handler.GetEmployees)
	http.HandleFunc("GET /employees/{id}", handler.GetEmployee)
	http.HandleFunc("PUT /employees/{id}", handler.UpdateEmployee)
	http.HandleFunc("DELETE /employees/{id}", handler.DeleteEmployee)

	// Ниже рабочая часть с MemoryStorage
	// memoryStorage := storages.NewMemoryStorage()
	// handler := NewHandler(memoryStorage)

	// http.HandleFunc("POST /employees/create", handler.CreateEmployee)
	// http.HandleFunc("GET /employees", handler.GetEmployees)
	// http.HandleFunc("GET /employees/{id}", handler.GetEmployee)
	// http.HandleFunc("PUT /employees/{id}", handler.UpdateEmployee)
	// http.HandleFunc("DELETE /employees/{id}", handler.DeleteEmployee)

	http.ListenAndServe(":8080", nil)
}
