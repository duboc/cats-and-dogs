package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/wavefronthq/go-metrics-wavefront/reporting"
	"github.com/wavefronthq/wavefront-sdk-go/application"
	"github.com/wavefronthq/wavefront-sdk-go/senders"

	"github.com/wavefronthq/wavefront-opentracing-sdk-go/reporter"
	"github.com/wavefronthq/wavefront-opentracing-sdk-go/tracer"

	nats "github.com/nats-io/nats.go"
)

var (
	nc *nats.Conn
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	var err error
	nc, err = nats.Connect(os.Getenv("NATS_ENDPOINT"))
	if err != nil {
		fmt.Println(err)
	}

	// Counter metrics registry
	metricRequestDurationDog := createDurationMetric("dog")
	metricRequestDurationCat := createDurationMetric("cat")

	//Metrics error
	counterErrorDog := createCounterErrorMetric("dog")
	counterErrorCat := createCounterErrorMetric("cat")

	// configure direct sender

	directCfg := &senders.DirectConfiguration{
		// Your Wavefront instance URL
		Server: os.Getenv("WAVEFRONT_INSTANCE"),

		// Wavefront API token created with direct ingestion permission
		Token: os.Getenv("WAVEFRONT_TOKEN"),

		// Optional: Override the batch size (in data points). Default: 10,000. Recommended not to exceed 40,000.
		BatchSize: 20000,

		// Optional: Override the max buffer size (in data points). Default: 50,000. Higher values could use more memory.
		MaxBufferSize: 50000,

		// Optional: Override the flush interval (in seconds). Default: 1 second
		FlushIntervalSeconds: 2,
	}

	// Create the direct sender
	sender, err := senders.NewDirectSender(directCfg)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	appTags := application.New("Petshop", "order")

	reporter := reporter.New(sender, appTags, reporter.Source(os.Getenv("WAVEFRONT_SOURCE")))

	tracer := tracer.New(reporter)

	opentracing.InitGlobalTracer(tracer)

	fmt.Println("Listen 9090 .....")
	http.HandleFunc("/api/dog", handler(metricRequestDurationDog, counterErrorDog))
	http.HandleFunc("/api/cat", handler(metricRequestDurationCat, counterErrorCat))
	http.ListenAndServe(":9090", nil)
}

func createDurationMetric(animal string) metrics.Timer {
	t := metrics.NewTimer()
	reporting.RegisterMetric(
		"request.duration", t, map[string]string{
			"animal": animal,
		})
	return t
}

func createCounterErrorMetric(animal string) metrics.Counter {
	c := metrics.NewCounter()
	reporting.RegisterMetric(
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
			if i := rand.Intn(5); i == 4 {
				w.WriteHeader(429)
				errorMetric.Inc(1)
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("Erro ao processar"))
				fmt.Println(nc.Publish("animals_errors", []byte("Erro ao processar")))
			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response{Animal: animal})
				str, _ := json.Marshal(response{Animal: animal})
				fmt.Println(nc.Publish("animals", []byte(str)))

			}
		})
	}

}

type response struct {
	Animal string `json:"animal"`
}
