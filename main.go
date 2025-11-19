package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Metrikler
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kobay_requests_total",
			Help: "Toplam HTTP istek sayisi",
		},
		[]string{"status"},
	)
	errorRateGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "kobay_error_rate_setting",
			Help: "Su anki hata orani ayari",
		},
	)
	currentErrorRate int
	mu               sync.RWMutex
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(errorRateGauge)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Ana Sayfa (Trafik)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		rate := currentErrorRate
		mu.RUnlock()

		if rand.Intn(100) < rate {
			requestCount.WithLabelValues("500").Inc()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Hata! Sistem bozuk (500).")
		} else {
			requestCount.WithLabelValues("200").Inc()
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Basarili. Sistem stabil (200).")
		}
	})

	// Kaos Butonu (/set-error/50 gibi)
	http.HandleFunc("/set-error/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 { return }
		rate, _ := strconv.Atoi(parts[2])
		
		mu.Lock()
		currentErrorRate = rate
		errorRateGauge.Set(float64(rate))
		mu.Unlock()
		fmt.Fprintf(w, "Hata orani %% %d yapildi.", rate)
	})

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Go Kobay Servisi 8080 portunda calisiyor...")
	http.ListenAndServe(":8080", nil)
}