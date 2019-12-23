package main

import (
	"os"
	"time"
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

var (
	topic = "my-topic"
	partition = 0
	dialer *kafka.Dialer
)



func main(){
	if len(os.Args) < 2 {
		fmt.Println("Favor passar r para leitura ou w para escrita com o valor a ser escrito Ex.: w 'mensagem'")
		os.Exit(1)
	}

	action := os.Args[1]

	dialer = &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	switch action {
	case "w": write(time.Now().String(), os.Args[2])
	case "r": read()
	}


}

func read(){
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		//GroupID:        "consumer-group-id",
		Topic:          "topic-A",
		Partition: 0,
		//MinBytes:  10e3, // 10KB
		MinBytes:  1e3, // 10KB
		MaxBytes:  10e6, // 10MB
		Dialer:         dialer,
	})
	//r.SetOffset(42)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	r.Close()
}

func write(key, value string){
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "topic-A",
		//Balancer: &kafka.Hash{},
		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
	})

	w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		},
	)

	w.Close()
}

