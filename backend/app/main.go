package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	//metrics "github.com/rcrowley/go-metrics"
	//wavefront "github.com/wavefronthq/go-metrics-wavefront"
	//"github.com/wavefronthq/wavefront-opentracing-sdk-go/reporter"
	//"github.com/wavefronthq/wavefront-opentracing-sdk-go/tracer"
	//"github.com/wavefronthq/wavefront-sdk-go/application"
	//"github.com/wavefronthq/wavefront-sdk-go/senders"
	
	metrics "github.com/rcrowley/go-metrics"
	"github.com/wavefronthq/go-metrics-wavefront/reporting"
	"github.com/wavefronthq/wavefront-sdk-go/application"
	"github.com/wavefronthq/wavefront-sdk-go/senders"

	"github.com/wavefronthq/wavefront-opentracing-sdk-go/reporter"
	"github.com/wavefronthq/wavefront-opentracing-sdk-go/tracer"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	// Counter metrics registry
	//metricRequestDurationDog := metrics.GetOrRegisterTimer("request.dog.duration", metrics.DefaultRegistry)
	//metricRequestDurationCat := metrics.GetOrRegisterTimer("request.cat.duration", metrics.DefaultRegistry)
	metricRequestDurationDog := createDurationMetric("dog")
	metricRequestDurationCat := createDurationMetric("cat")

	//Metrics error
	//counterErrorDog := metrics.GetOrRegisterCounter("errorDog", metrics.DefaultRegistry)
	//counterErrorCat := metrics.GetOrRegisterCounter("errorCat", metrics.DefaultRegistry)
	counterErrorDog := createCounterErrorMetric("dog")
	counterErrorCat := createCounterErrorMetric("cat")

	//hostTags := map[string]string{
	//	"source": "petshop",
	//}

	// configure direct sender

	directCfg := &senders.DirectConfiguration{
		// Your Wavefront instance URL
		Server: "https://vmware.wavefront.com",

		// Wavefront API token created with direct ingestion permission
		Token: "b573b1b1-d42a-4518-9611-5a5668891671",

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
	}
	// Use the direct sender to send data

	//app metrics https://github.com/wavefrontHQ/wavefront-opentracing-sdk-go

	appTags := application.New("Petshop", "order")

	reporter := reporter.New(sender, appTags, reporter.Source("petshop.app.com"))

	tracer := tracer.New(reporter)

	opentracing.InitGlobalTracer(tracer)

	////server := os.Getenv("WAVEFRONT_INSTANCE")
	////token := os.Getenv("WAVEFRONT_TOKEN")
	////go wavefront.WavefrontDirect(metrics.DefaultRegistry, 5*time.Second, hostTags, "zoologico", server, token)
	//// report to a Wavefront proxy
	//addr, _ := net.ResolveTCPAddr("tcp", os.Getenv(""))
	//go wavefront.WavefrontProxy(metrics.DefaultRegistry, 5*time.Second, hostTags, "petshop", addr)
	////go metrics.Log(metrics.DefaultRegistry, 10*time.Second, log.New(os.Stdout, "metrics: ", log.Lmicroseconds))

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
			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response{Animal: animal})
			}
		})
	}

}

type response struct {
	Animal string `json:"animal"`
}
