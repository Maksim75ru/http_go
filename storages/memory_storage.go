package storages

import (
	"errors"
	"strconv"
	"sync"
)

type MemoryStorage struct {
	counter int
	data    map[string]Employee
	sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data:    make(map[string]Employee),
		counter: 1,
	}
}

func (s *MemoryStorage) Create(e Employee) Employee {
	s.Lock()
	defer s.Unlock()
	e.Id = strconv.Itoa(s.counter)
	s.data[e.Id] = e
	s.counter++
	return s.data[e.Id]
}

func (s *MemoryStorage) Get(id string) (Employee, error) {
	employee, exists := s.data[id]
	if !exists {
		return Employee{}, errors.New("Employee with such Id doesn't exist")
	}
	return employee, nil
}

func (s *MemoryStorage) GetAll() []Employee {
	employees := []Employee{}
	for _, employee := range s.data {
		employees = append(employees, employee)
	}
	return employees
}

func (s *MemoryStorage) Update(id string, e Employee) (bool, error) {
	s.Lock()
	_, exists := s.data[id]
	if !exists {
		return false, errors.New("Employee with such Id doesn't exist")
	}
	s.data[id] = e
	s.Unlock()
	return true, nil
}

func (s *MemoryStorage) Delete(id string) (bool, error) {
	s.Lock()
	_, exists := s.data[id]
	if !exists {
		return false, errors.New("Employee with such Id doesn't exist")
	}
	delete(s.data, id)
	s.Unlock()
	return true, nil
}
