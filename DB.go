package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	DB *mongo.Client
}

type VectorClkDB struct {
	port  uint64
	Key   string
	Value string
	VC    [N]int
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

func (DB Database) Get(port uint64, key string) (string, bool) {
	// collection := DB.DB.Database("kv").Collection("kvv")
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	val := DB.getVersions(port, key)

	// fmt.Println("Returned values:", values)

	// result := make(map[string]interface{})
	// err := collection.FindOne(ctx, bson.D{{"key", key}}).Decode(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// value := result["value"].(string)
	// log.Printf("Get value [%s] \n", value)

	return val, true
}

func (DB Database) Put(port uint64, key string, value string) {
	collection := DB.DB.Database("kv").Collection("kvv")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// build VectorClkDB
	latest_vc := DB.GetVC(port, key)
	latest_vc[port%8081] += 1

	// initialize vector clock
	doc := bson.D{{"port", port}, {"key", key}, {"value", value}, {"vc", latest_vc}}

	putResult, _ := collection.InsertOne(ctx, doc)
	log.Printf("Inserted key [%s] value [%s], insert result: [%s] \n", key, value, putResult)
}

func (DB Database) GetVC(port uint64, key string) [N]int {
	collection := DB.DB.Database("kv").Collection("kvv")
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	filter := bson.D{{"port", port}, {"key", key}}
	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		panic(err)
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	output := [N]int{0, 0, 0, 0, 0}

	if len(results) == 0 {
		return output
	}

	for _, result := range results {
		doc, err := bson.Marshal(result)
		if err != nil {
			panic(err)
		}
		var temp VectorClkDB
		err = bson.Unmarshal(doc, &temp)
		if err != nil {
			panic(err)
		}
		curr_vc := temp.VC

		if CompareVCs(output, curr_vc) == 2 {
			output = curr_vc
		}
	}
	return output
}

func (DB Database) getVersions(port uint64, key string) string {
	collection := DB.DB.Database("kv").Collection("kvv")
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	filter := bson.D{{"port", port}, {"key", key}}
	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		panic(err)
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	if len(results) == 0 {
		return ""
	}

	output := [N]int{0, 0, 0, 0, 0}
	value := ""
	for _, result := range results {
		doc, err := bson.Marshal(result)
		if err != nil {
			panic(err)
		}
		var temp VectorClkDB
		err = bson.Unmarshal(doc, &temp)
		if err != nil {
			panic(err)
		}
		curr_vc := temp.VC

		if CompareVCs(output, curr_vc) == 2 {
			output = curr_vc
			value = temp.Value
		}
	}

	fmt.Print("Conflict Vector Clock: ", output)

	return value
}

func CompareVCs(arr1 [N]int, arr2 [N]int) int {
	// return values
	// 1: arr1 is greater
	// 2: arr2 is greater
	// 3: conflict

	arr1_greater := true
	for i := 0; i < N; i++ {
		if !(arr1[i] >= arr2[i]) {
			arr1_greater = false
			break
		}
	}

	arr2_greater := true
	for i := 0; i < N; i++ {
		if !(arr2[i] >= arr1[i]) {
			arr2_greater = false
			break
		}
	}

	if !arr1_greater && !arr2_greater {
		return 3
	}
	if arr1_greater {
		return 1
	}
	if arr2_greater {
		return 2
	}
	return 0
}
