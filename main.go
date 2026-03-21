package main

import (
	"fmt"
	"go-app/controllers"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)


var LogPath = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: %s (%s)", r.Host, r.RequestURI, r.Method))
		next.ServeHTTP(w, r)
	})
}

var (
    // Счетчик запросов
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"}, // лейблы: метод, путь, статус
    )
    
    // Гистограмма для длительности запросов
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}, // корзины для латентности
        },
        []string{"method", "path"},
    )
)

// Middleware для сбора метрик
func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Создаем обертку для ResponseWriter, чтобы перехватить статус код
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        // Передаем управление дальше
        next.ServeHTTP(wrapped, r)
        
        // После обработки запроса собираем метрики
        duration := time.Since(start).Seconds()
        
        // Путь для метрик (нужно нормализовать, чтобы не было /users/123)
        path := normalizePath(r.URL.Path)
        
        // Увеличиваем счетчик запросов
        httpRequestsTotal.WithLabelValues(r.Method, path, fmt.Sprint(wrapped.statusCode)).Inc()
        
        // Добавляем гистограмму
        httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
    })
}

// Обертка для ResponseWriter, чтобы перехватить статус код
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// Функция нормализации путей
// Вместо /users/123 -> /users/:id
func normalizePath(path string) string {
    // Это простейший вариант. В реальном проекте используйте роутер с паттернами
    // Пример для чистого net/http
    if len(path) > 0 && path[len(path)-1] == '/' {
        path = path[:len(path)-1]
    }
    return path
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/notes", controllers.NoteQuery).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/notes", controllers.NoteCreate).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/notes/{id}", controllers.NoteRetrieve).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/notes/{id}", controllers.NoteUpdate).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/notes/{id}", controllers.NoteDelete).Methods(http.MethodDelete, http.MethodOptions)
	
	router.Handle("/metrics", promhttp.Handler())
	
	router.Use(LogPath)


	handler := metricsMiddleware(router)

	log.Info("Listening on 8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
