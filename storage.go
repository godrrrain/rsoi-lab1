package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID      int    `bson:"_id"`
	Name    string `bson:"name" binding:"required"`
	Age     int    `bson:"age" binding:"required"`
	Address string `bson:"address" binding:"required"`
	Work    string `bson:"work" binding:"required"`
}

type Storage interface {
	Insert(e *Person)
	Get(id int) (Person, error)
	Update(e *Person) error
	Delete(id int) error
	GetAll() []Person
}

type MongoDbStorage struct {
	c *mongo.Collection
	sync.Mutex
}

func NewMongoDbStorage() *MongoDbStorage {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var mongoUri = os.Getenv("MONGOURI")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Fatalf("Error ping to MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB!")

	collection := client.Database("lab").Collection("persons")

	return &MongoDbStorage{
		c: collection,
	}
}

func NewTestMongoDbStorage() *MongoDbStorage {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var mongoUri = os.Getenv("MONGOURI")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	collection := client.Database("testing").Collection("persons")

	return &MongoDbStorage{
		c: collection,
	}
}

func (s *MongoDbStorage) Get(id int) (Person, error) {
	s.Lock()
	defer s.Unlock()

	var person Person

	filter := bson.D{{Key: "_id", Value: id}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.c.FindOne(ctx, filter).Decode(&person)
	if err == mongo.ErrNoDocuments {
		return person, errors.New("person not found")
	} else if err != nil {
		return person, err
	}

	return person, nil
}

func (s *MongoDbStorage) Insert(e *Person) {
	s.Lock()

	count, _ := s.c.CountDocuments(context.Background(), bson.D{})
	e.ID = int(count)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, _ := s.c.InsertOne(ctx, e)
	id := res.InsertedID
	fmt.Printf("Id = %v", id)

	s.Unlock()
}

func (s *MongoDbStorage) GetAll() []Person {
	s.Lock()
	defer s.Unlock()

	count, _ := s.c.CountDocuments(context.Background(), bson.D{})
	persons := make([]Person, count)

	cur, err := s.c.Find(context.Background(), bson.D{})
	if err != nil {
		//logging error
	}
	defer cur.Close(context.Background())
	number := 0
	for cur.Next(context.Background()) {
		var result Person
		err := cur.Decode(&result)
		if err != nil {
			//logging error
		}
		persons[number] = result
		number = number + 1
	}
	if err := cur.Err(); err != nil {
		//logging error
	}

	return persons
}

func (s *MongoDbStorage) Update(e *Person) error {
	s.Lock()
	defer s.Unlock()

	filter := bson.D{{Key: "_id", Value: e.ID}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := s.c.ReplaceOne(ctx, filter, e)
	fmt.Println(result)
	fmt.Println(err)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("person not found")
	}

	return nil
}

func (s *MongoDbStorage) Delete(id int) error {
	s.Lock()
	defer s.Unlock()

	filter := bson.D{{Key: "_id", Value: id}}
	result, err := s.c.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("person not found")
	}

	return nil
}
