package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"http_go/storages"
	"net/http"
	"strings"
)

type Handler struct {
	storage storages.Storage
}

func NewHandler(storage storages.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	someStorage := h.storage

	var e storages.Employee
	if err := dec.Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newEmployee := someStorage.Create(e)

	fmt.Printf("You created new employee with ID=%s\n", newEmployee.Id)

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

	someStorage := h.storage
	employee, err := someStorage.Get(string(eId))
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

func (h *Handler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	eId, err := getEmployeeId(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var e storages.Employee
	if err := dec.Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	memoryStorage := h.storage
	isUpdated, err := memoryStorage.Update(eId, e)
	if !isUpdated {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("You updated employee with ID=%s\n", eId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created

	if err := json.NewEncoder(w).Encode(e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getEmployeeId(urlPath string) (string, error) {
	idStr := strings.TrimPrefix(urlPath, "/employees/")

	if idStr == "" || idStr == urlPath {
		return "", errors.New("endpoint doesnt contain ID")
	}

	return idStr, nil
}
