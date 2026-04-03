package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// Middleware для логирования
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("[%s] %s %s", r.Method, r.URL.Path, time.Since(start))
        next(w, r)
    }
}

// Обработчик для проверки
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from Go with logging!")
}

func main() {
    http.HandleFunc("/", loggingMiddleware(helloHandler))
    fmt.Println("Go сервер запущен на порту 8081")
    log.Fatal(http.ListenAndServe(":8081", nil))
}