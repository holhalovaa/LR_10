// Package main API Gateway для маршрутизации запросов
//
// @title           API Gateway
// @version         1.0
// @description     API Gateway for routing to Go and Python services
// @host            localhost:8082
// @BasePath        /
//
// @securityDefinitions.basic  BasicAuth
package main

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go_service/docs"
)

// healthHandler godoc
// @Summary      Health check
// @Description  Returns the status of the API Gateway
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

// goProxy godoc
// @Summary      Proxy to Go service
// @Description  Forwards requests to the Go service on port 8081
// @Success      200  {string}  string  "Response from Go service"
// @Router       /go/{path} [get]
func goProxy(c *gin.Context) {
	goURL, _ := url.Parse("http://localhost:8081")
	proxy := httputil.NewSingleHostReverseProxy(goURL)
	c.Request.URL.Path = c.Param("proxyPath")
	if c.Request.URL.Path == "" {
		c.Request.URL.Path = "/"
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

// pythonProxy godoc
// @Summary      Proxy to Python service
// @Description  Forwards requests to the Python FastAPI service on port 8000
// @Success      200  {string}  string  "Response from Python service"
// @Router       /python/{path} [get]
func pythonProxy(c *gin.Context) {
	pythonURL, _ := url.Parse("http://localhost:8000")
	proxy := httputil.NewSingleHostReverseProxy(pythonURL)
	c.Request.URL.Path = c.Param("proxyPath")
	if c.Request.URL.Path == "" {
		c.Request.URL.Path = "/"
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {
	r := gin.Default()

	r.GET("/health", healthHandler)
	r.Any("/go/*proxyPath", goProxy)
	r.Any("/python/*proxyPath", pythonProxy)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("API Gateway запущен на порту 8082")
	log.Fatal(r.Run(":8082"))
}
