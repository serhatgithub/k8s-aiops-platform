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
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "demo_requests_total", // DÜZELTİLDİ: Artık 'demo' ismiyle yayın yapıyor
			Help: "Toplam istek sayisi",
		},
		[]string{"status"},
	)
	errorRateGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "demo_error_rate",
			Help: "Hata orani ayari",
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
	// Ana Sayfa Handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		rate := currentErrorRate
		mu.RUnlock()

		if rand.Intn(100) < rate {
			requestCount.WithLabelValues("500").Inc()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Hata! Sistem patladi (500).")
		} else {
			requestCount.WithLabelValues("200").Inc()
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Stabil. Her sey yolunda (200).")
		}
	})

	// Hata Ayarlama Handler
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
	fmt.Println("Go Demo App 8080 portunda calisiyor...")
	http.ListenAndServe(":8080", nil)
}