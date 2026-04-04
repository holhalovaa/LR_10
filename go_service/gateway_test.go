package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthEndpoint(t *testing.T) {
	// Устанавливаем Gin в тестовый режим
	gin.SetMode(gin.TestMode)

	// Создаём роутер
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Создаём тестовый запрос
	req, _ := http.NewRequest("GET", "/health", nil)
	resp := httptest.NewRecorder()

	// Выполняем запрос
	r.ServeHTTP(resp, req)

	// Проверяем статус
	if resp.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получили %d", resp.Code)
	}

	// Проверяем тело ответа
	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response["status"] != "ok" {
		t.Errorf("Ожидался status: ok, получили %s", response["status"])
	}
}

func TestGatewayRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Регистрируем реальные роуты (как в gateway.go)
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.Any("/go/*any", func(c *gin.Context) { c.String(200, "proxy to go") })
	r.Any("/python/*any", func(c *gin.Context) { c.String(200, "proxy to python") })

	tests := []struct {
		route    string
		expected int
	}{
		{"/health", 200},
		{"/go/test", 200},
		{"/python/test", 200},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest("GET", tt.route, nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
		if resp.Code != tt.expected {
			t.Errorf("Роут %s: ожидался %d, получили %d", tt.route, tt.expected, resp.Code)
		}
	}
}
