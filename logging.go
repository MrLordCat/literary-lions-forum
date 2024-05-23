package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Middleware для логирования ответов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логгер, который записывает данные перед отправкой клиенту
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		next.ServeHTTP(lw, r)

		// Логируем статус и данные после выполнения обработчика
		log.Printf("[%s] %s %d", r.Method, r.URL.Path, lw.statusCode)
		if lw.data != nil {
			logData, _ := json.MarshalIndent(lw.data, "", "  ")
			log.Printf("Data: %s", logData)
		}
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	data       interface{}
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	// Здесь можно десериализовать JSON для лучшего вывода, если ответ в JSON
	if json.Unmarshal(data, &lrw.data) != nil {
		lrw.data = nil // Если не удалось десериализовать, игнорируем
	}
	return lrw.ResponseWriter.Write(data)
}
