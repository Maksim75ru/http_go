package main

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Employee struct {
	Id         int
	Name       string
	Sex        string
	Age        int
	Salary     int
	Department string
}

type Storage interface {
	Create(e Employee) Employee
	Get(id int) (Employee, error)
	GetAll() []Employee
	Update(id int, e Employee) (bool, error)
	Delete(id int) (bool, error)
}

type MemoryStorage struct {
	counter int
	data    map[int]Employee
	sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data:    make(map[int]Employee),
		counter: 1,
	}
}

func (s *MemoryStorage) Create(e Employee) Employee {
	s.Lock()
	e.Id = s.counter
	s.data[e.Id] = e
	s.counter++
	s.Unlock()
	return s.data[e.Id]
}

func (s *MemoryStorage) Get(id int) (Employee, error) {
	e, exists := s.data[id]
	if !exists {
		return Employee{}, errors.New("Employee with such Id doesn't exist")
	}
	return e, nil
}

func (s *MemoryStorage) GetAll() []Employee {
	employees := []Employee{}
	for _, employee := range s.data {
		employees = append(employees, employee)
	}
	return employees
}

func (s *MemoryStorage) Update(id int, e Employee) (bool, error) {
	s.Lock()
	_, exists := s.data[id]
	if !exists {
		return false, errors.New("Employee with such Id doesn't exist")
	}
	s.data[id] = e
	s.Unlock()
	return true, nil
}

func (s *MemoryStorage) Delete(id int) (bool, error) {
	s.Lock()
	_, exists := s.data[id]
	if !exists {
		return false, errors.New("Employee with such Id doesn't exist")
	}
	delete(s.data, id)
	s.Unlock()
	return true, nil
}

type MongoDbStorage struct {
	coll interface{}
}

func NewMongoDbStorage(dbName, collection, mongoDBURI string) *MongoDbStorage {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return nil
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. ")
		return nil
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	coll := client.Database("RogaAndKopita").Collection("employees")

	return &MongoDbStorage{coll: coll}
}
