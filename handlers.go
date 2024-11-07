package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	memoryStorage := NewMemoryStorage()

	var e Employee
	if err := dec.Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newEmployee := memoryStorage.Create(e)

	fmt.Printf("You created new employee with ID=%d\n", newEmployee.Id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created

	if err := json.NewEncoder(w).Encode(newEmployee); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
