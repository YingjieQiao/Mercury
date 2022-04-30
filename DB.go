package main

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Database struct {
	DB *mongo.Client
}

func NewDatabase() *Database {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("No mongodb uri")
	}

	mongoDB, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	store := &Database{DB: mongoDB}
	return store
}

func (DB Database) Get(key string) (string, bool) {
	collection := DB.DB.Database("kv").Collection("kv")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := make(map[string]interface{})
	err := collection.FindOne(ctx, bson.D{{"key", key}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	value := result["value"].(string)
	log.Printf("Get value [%s] \n", value)

	return value, true
}

func (DB Database) Put(key string, value string) {
	collection := DB.DB.Database("kv").Collection("kv")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	doc := bson.D{{"key", key}, {"value", value}}

	putResult, _ := collection.InsertOne(ctx, doc)
	log.Printf("Inserted key [%s] value [%s], insert result: [%s] \n", key, value, putResult)
}
