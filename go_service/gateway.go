package main

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Middleware логирования
	r.Use(func(c *gin.Context) {
		log.Printf("Gateway: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// Прокси на Python сервис (FastAPI на порту 8000)
	pythonURL, _ := url.Parse("http://localhost:8000")
	pythonProxy := httputil.NewSingleHostReverseProxy(pythonURL)

	// Все запросы, начинающиеся с /python/, перенаправляем в FastAPI
	r.Any("/python/*proxyPath", func(c *gin.Context) {
		// Убираем /python из пути, чтобы FastAPI получил правильный маршрут
		c.Request.URL.Path = c.Param("proxyPath")
		if c.Request.URL.Path == "" {
			c.Request.URL.Path = "/"
		}
		pythonProxy.ServeHTTP(c.Writer, c.Request)
	})

	// Прокси на Go сервис (порт 8081)
	goURL, _ := url.Parse("http://localhost:8081")
	goProxy := httputil.NewSingleHostReverseProxy(goURL)

	// Все запросы, начинающиеся с /go/, перенаправляем в Go сервис
	r.Any("/go/*proxyPath", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("proxyPath")
		if c.Request.URL.Path == "" {
			c.Request.URL.Path = "/"
		}
		goProxy.ServeHTTP(c.Writer, c.Request)
	})

	// Проверка здоровья Gateway
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	fmt.Println("API Gateway запущен на порту 8082")
	log.Fatal(r.Run(":8082"))
}
