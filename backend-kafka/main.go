package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	wavefront "github.com/wavefronthq/go-metrics-wavefront"

	"github.com/segmentio/kafka-go"
)

var (
	dialer = &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
	kafkaWriter *kafka.Writer
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("KAFKA_URL")},
		Topic:   "topic-A",
		//Balancer: &kafka.Hash{},
		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
	})

	// Counter metrics registry
	metricRequestDurationDog := createDurationMetric("dog")
	metricRequestDurationCat := createDurationMetric("cat")

	//Metrics error
	counterErrorDog := createCounterErrorMetric("dog")
	counterErrorCat := createCounterErrorMetric("cat")

	hostTags := map[string]string{
		"source": "petshop",
	}

	// report to a Wavefront proxy
	addr, _ := net.ResolveTCPAddr("tcp", os.Getenv("WF_PROXY"))
	go wavefront.WavefrontProxy(metrics.DefaultRegistry, 5*time.Second, hostTags, "petshop", addr)

	fmt.Println("Listen 9090 .....")
	http.HandleFunc("/api/dog", handler(metricRequestDurationDog, counterErrorDog))
	http.HandleFunc("/api/cat", handler(metricRequestDurationCat, counterErrorCat))
	http.ListenAndServe(":9090", nil)
	kafkaWriter.Close()
}

func createDurationMetric(animal string) metrics.Timer {
	t := metrics.NewTimer()
	wavefront.RegisterMetric(
		"request.duration", t, map[string]string{
			"animal": animal,
		})
	return t
}

func createCounterErrorMetric(animal string) metrics.Counter {
	c := metrics.NewCounter()
	wavefront.RegisterMetric(
		"error", c, map[string]string{
			"animal": animal,
		})
	return c
}

//handler Devolve uma funcao que vai tratar a requisicao e gerar as metricas
// Recebe como parametro um
//   durationMetric: Onde serao geradas as metricas referente a duracao da requisicao
//   errorMetric: Onde serao contados os erros na requisicao
func handler(durationMetric metrics.Timer, errorMetric metrics.Counter) func(http.ResponseWriter, *http.Request) {
	// Importante: Aqui sera criado uma FUNCAO que sera retornada por este metodo
	return func(w http.ResponseWriter, r *http.Request) {
		//Envolve toda a logica da request no durationMetric, para tirar as metricas de duracao
		durationMetric.Time(func() {
			log.Printf("Received request %s\n", r.URL.String())
			animal := r.URL.Path[len("/api/"):]
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			res := response{Time: time.Now().String(), Animal: animal, ResultCode: 200}
			if i := rand.Intn(5); i == 4 {
				res.ResultCode = 429
				w.WriteHeader(429)
				errorMetric.Inc(1)
			}
			arr, _ := json.Marshal(res)
			sendToKafka(r.Context(), "request", arr)
			w.Header().Set("Content-Type", "application/json")
			w.Write(arr)
		})
	}

}

func sendToKafka(ctx context.Context, key string, value []byte){
	kafkaWriter.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: value,
		},
	)

}

type response struct {
	Time string `json:"time"`
	Animal string `json:"animal"`
	ResultCode int `json:"result_code"`
}
