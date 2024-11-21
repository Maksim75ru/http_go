package storages

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbStorage struct {
	coll *mongo.Collection
}

func NewMongoDbStorage(dbName, collection string) *MongoDbStorage {
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

func (s *MongoDbStorage) Create(e Employee) Employee {
	// Генерация нового ObjectID
	e.Id = primitive.NewObjectID().Hex()

	res, err := s.coll.InsertOne(context.TODO(), e)
	if err != nil {
		fmt.Println("Error inserting employee:", err)
		return Employee{} // Вернуть пустого сотрудника в случае ошибки
	}

	// Извлечение ID вставленного документа
	insertedID := res.InsertedID
	fmt.Printf("Inserted document ID: %v\n", insertedID)

	eID, ok := insertedID.(primitive.ObjectID) // Здесь используем bson.ObjectID
	if !ok {
		fmt.Println("Some problem with converted bson.ObjectID")
		return Employee{} // Вернуть пустого сотрудника в случае ошибки
	}

	e.Id = eID.Hex()

	return e
}

func (s *MongoDbStorage) Get(id string) (Employee, error) {
	var employee Employee

	filter := bson.M{"id": id}
	err := s.coll.FindOne(context.TODO(), filter).Decode(&employee)

	if err != nil {
		return Employee{}, err // Возвращаем ошибку, если не удалось найти
	}

	return employee, nil
}

func (s *MongoDbStorage) GetAll() []Employee {
	employees := []Employee{}

	filter := bson.M{}
	cursor, err := s.coll.Find(context.TODO(), filter)

	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO()) // Закрыть курсор после завершения работы

	if err := cursor.All(context.TODO(), &employees); err != nil {
		panic(err)
	}

	return employees
}

func (s *MongoDbStorage) Update(id string, e Employee) (bool, error) {
	_, exists := s.Get(id)
	if exists != nil {
		return false, errors.New("Employee with such Id doesn't exist")
	}

	filter := bson.M{"id": id}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: e.Name},
				{Key: "sex", Value: e.Sex},
				{Key: "department", Value: e.Department},
				{Key: "age", Value: e.Age},
				{Key: "salary", Value: e.Salary},
			},
		},
	}

	_, err := s.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("Error updating employee:", err)
		return false, err
	}

	return true, nil
}

func (s *MongoDbStorage) Delete(id string) (bool, error) {
	_, exists := s.Get(id)
	if exists != nil {
		return false, errors.New("Employee with such Id doesn't exist")
	}

	filter := bson.M{"id": id}
	_, err := s.coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error deleting employee:", err)
		return false, err
	}

	return true, nil
}
