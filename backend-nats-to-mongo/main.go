package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	nc *nats.Conn
)

func main() {
	var err error
	nc, err = nats.Connect(os.Getenv("NATS_ENDPOINT"))
	if err != nil {
		fmt.Println(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_ENDPOINT")))
	collection := client.Database("testing").Collection("animals")

	nc.Subscribe("animals", func(m *nats.Msg) {
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		bdoc := bson.D{}
		bson.UnmarshalExtJSON(m.Data, true, &bdoc)
		_, err := collection.InsertOne(ctx, bdoc)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Message sent to mongo %s\n", string(m.Data))
		}

	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
