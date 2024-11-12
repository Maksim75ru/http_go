package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	uri := os.Getenv("MONGODB_URI")
	mongoDbStorage := NewMongoDbStorage("RogaAndKopita", "employees", uri)
	// handler := NewHandler(mongoDbStorage)
	fmt.Println(mongoDbStorage)

	// http.HandleFunc("POST /employee/create", handler.CreateEmployee)

	// Ниже рабочая часть с MemoryStorage
	// memoryStorage := NewMemoryStorage()
	// handler := NewHandler(memoryStorage)

	// http.HandleFunc("POST /employee/create", handler.CreateEmployee)
	// http.HandleFunc("GET /employees", handler.GetEmployees)
	// http.HandleFunc("GET /employee/{id}", handler.GetEmployee)
	// http.HandleFunc("PUT /employee/{id}", handler.UpdateEmployee)
	// http.HandleFunc("DELETE /employee/{id}", handler.DeleteEmployee)

	http.ListenAndServe(":8080", nil)
}
