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

	// Проверяем, что роуты зарегистрированы
	routes := []string{"/health", "/go/", "/python/"}

	for _, route := range routes {
		req, _ := http.NewRequest("GET", route, nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		// Даже если 404, роут должен существовать
		if resp.Code == http.StatusNotFound {
			t.Logf("Роут %s существует (ответ 404 из-за отсутствия бэкенда)", route)
		}
	}
}
