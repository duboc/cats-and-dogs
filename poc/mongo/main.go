package main

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/bson"
	"time"
	"context"
	"fmt"
	"os"
	"log"

)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Favor passar r <filtro> para leitura ou w <mensagem> para escrita")
		os.Exit(1)
	}

	action := os.Args[1]
	param := os.Args[2]

	client, err := mongo.NewClient("mongodb://localhost:27017")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil { log.Fatal(err) }

	collection := client.Database("testing").Collection("numbers")

	switch action {
	case "w":
		insert(collection, bson.M{"timestamp": time.Now().String(), "name": param})
	case "r": read(collection,bson.M{"name": param})
	}

}

func insert(collection *mongo.Collection, msg bson.M) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, msg)
	if err != nil { log.Fatal(err) }
	id := res.InsertedID
	fmt.Println("Inserted",id)
}

func read(collection *mongo.Collection, filter bson.M) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, filter)
	if err != nil { log.Fatal(err) }
	defer cur.Close(ctx)
	for cur.Next(ctx) {
	   var result bson.M
	   err := cur.Decode(&result)

	   if err != nil { log.Fatal(err) }
	   fmt.Println("Teste: ", result)
	}
	if err := cur.Err(); err != nil {
	  log.Fatal(err)
	}
}
