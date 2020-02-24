package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	kafka "github.com/segmentio/kafka-go"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var (
	dialer     *kafka.Dialer
	collection *mongo.Collection
)

func main() {

	dialer = &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
	//client, err := mongo.NewClient("mongodb://localhost:27017")
	client, err := mongo.NewClient(os.Getenv("MONGO_URL"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("catdog").Collection("animals")

	kafkaRead()

}

func kafkaRead() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_URL")},
		//GroupID:        "consumer-group-id",
		Topic:     "topic-A",
		Partition: 0,
		//MinBytes:  10e3, // 10KB
		MinBytes: 1e3,  // 10KB
		MaxBytes: 10e6, // 10MB
		Dialer:   dialer,
	})
	//r.SetOffset(42)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		var animal Animal
		err = json.Unmarshal(m.Value, &animal)
		if err != nil {
			log.Printf("ERROR: Error on process message '%s' '%s'", string(m.Value), err)
			continue
		}
		err = mongoInsert(collection, animal)
		if err != nil {
			log.Printf("ERROR: Error on insert in mongo message '%s' '%s'", string(m.Value), err)
		}
	}

	r.Close()
}

func mongoInsert(collection *mongo.Collection, msg interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	b, err := bson.Marshal(msg)
	if err != nil {
		return err
	}
	err = mongoInsert(collection, b)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

type Animal struct {
	Time       string `json:"time"`
	Animal     string `json:"animal"`
	ResultCode int    `json:"result_code"`
}
