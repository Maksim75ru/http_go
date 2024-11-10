package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

	memoryStorage := h.storage

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
		return
	}
}

func (h *Handler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	eId, err := getEmployeeId(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	memoryStorage := h.storage
	employee, err := memoryStorage.Get(eId)
	if err != nil {
		http.NotFound(w, r) // 404 Not Found
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employee); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	employees := h.storage.GetAll()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	eId, err := getEmployeeId(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	memoryStorage := h.storage
	isDeleted, err := memoryStorage.Delete(eId)
	if err != nil {
		http.NotFound(w, r) // 404 Not Found
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(isDeleted); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getEmployeeId(urlPath string) (int, error) {
	idStr := strings.TrimPrefix(urlPath, "/employee/")
	if idStr == "" || idStr == urlPath {
		return 0, errors.New("endpoint doesnt contain ID")
	}

	eId, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("invalid employee ID")
	}
	return eId, nil
}
