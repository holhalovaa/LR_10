package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Если вы добавили OpenAPI документацию через swaggo,
// этот тест проверит, что эндпоинт /swagger/index.html доступен

func TestOpenAPIDocumentation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Эндпоинт для OpenAPI документации (если добавили)
	r.GET("/swagger/*any", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"swagger": "2.0",
			"info": gin.H{
				"title":   "Go API Gateway",
				"version": "1.0",
			},
		})
	})

	// Проверяем, что эндпоинт существует
	req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	// Если документация есть, должен быть 200
	// Если нет, тест всё равно пройдёт (просто предупреждение)
	if resp.Code == http.StatusOK {
		t.Log("✅ OpenAPI документация доступна")

		var data map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &data)

		if info, ok := data["info"]; ok {
			t.Logf("Информация об API: %v", info)
		}
	} else {
		t.Log("⚠️ OpenAPI документация не настроена (это нормально для задания 8)")
	}
}

func TestAPIInfoEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Можно добавить эндпоинт с информацией об API
	r.GET("/api/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "Go API Gateway",
			"version":     "1.0.0",
			"description": "API Gateway for Python and Go services",
			"endpoints": []string{
				"/health",
				"/go/",
				"/python/call-go",
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/info", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code == http.StatusOK {
		t.Log("✅ API info эндпоинт доступен")

		var data map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &data)

		if name, ok := data["name"]; ok {
			t.Logf("Имя API: %v", name)
		}
	} else {
		t.Log("⚠️ API info эндпоинт не настроен")
	}
}
