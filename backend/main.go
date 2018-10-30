package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	wavefront "github.com/wavefronthq/go-metrics-wavefront"
)

var (
	metricRequestDuration metrics.Timer
	sleepTime             time.Duration
)

func main() {
	if s, ok := os.LookupEnv("SLEEP_SECONDS"); ok {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Error on convert SLEEP_SECONDS [%s] [%+v]", s, err)
		}
		sleepTime = time.Duration(i) * time.Second
	} else {
		sleepTime = 10
	}
	log.Printf("Sleeptime [%v]", sleepTime)
	metricRequestDuration = metrics.GetOrRegisterTimer("request.duration", metrics.DefaultRegistry)

	hostTags := map[string]string{
		"source": "go-metrics-test",
	}
	server := "https://vmware.wavefront.com"
	token := "b573b1b1-d42a-4518-9611-5a5668891671"
	go wavefront.WavefrontDirect(metrics.DefaultRegistry, 5*time.Second, hostTags, "direct.prefix", server, token)
	go metrics.Log(metrics.DefaultRegistry, 60*time.Second, log.New(os.Stdout, "metrics: ", log.Lmicroseconds))

	fmt.Println("Listen 9090 .....")
	http.HandleFunc("/api/", handlerWrapper)
	http.ListenAndServe(":9090", nil)
}

func handlerWrapper(w http.ResponseWriter, r *http.Request) {
	metricRequestDuration.Time(func() {
		handler(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request %s\n", r.URL.String())
	time.Sleep(10 * time.Second)
	animal := r.URL.Path[len("/api/"):]
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response{Animal: animal})
}

type response struct {
	Animal string `json:"animal"`
}
