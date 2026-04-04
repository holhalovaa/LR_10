package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestSwaggerUIAvailable(t *testing.T) {
	// Запускаем сервер в горутине
	go main()
	time.Sleep(2 * time.Second)

	// Проверяем, что страница Swagger UI доступна
	resp, err := http.Get("http://localhost:8082/swagger/index.html")
	if err != nil {
		t.Fatalf("Не удалось подключиться к Swagger UI: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Swagger UI вернул статус %d, ожидался 200", resp.StatusCode)
	}
}

func TestSwaggerJSONExists(t *testing.T) {
	// Запускаем сервер в горутине
	go main()
	time.Sleep(2 * time.Second)

	// Проверяем, что doc.json доступен
	resp, err := http.Get("http://localhost:8082/swagger/doc.json")
	if err != nil {
		t.Fatalf("Не удалось подключиться к doc.json: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("doc.json вернул статус %d, ожидался 200", resp.StatusCode)
	}

	// Проверяем, что ответ — валидный JSON
	var swaggerDoc map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&swaggerDoc); err != nil {
		t.Errorf("doc.json не является валидным JSON: %v", err)
	}

	// Проверяем, что в документации есть основные поля
	if _, ok := swaggerDoc["paths"]; !ok {
		t.Error("В doc.json отсутствует поле 'paths'")
	}

	if _, ok := swaggerDoc["info"]; !ok {
		t.Error("В doc.json отсутствует поле 'info'")
	}
}

func TestSwaggerHasEndpoints(t *testing.T) {
	// Запускаем сервер в горутине
	go main()
	time.Sleep(2 * time.Second)

	resp, err := http.Get("http://localhost:8082/swagger/doc.json")
	if err != nil {
		t.Fatalf("Не удалось подключиться к doc.json: %v", err)
	}
	defer resp.Body.Close()

	var swaggerDoc map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&swaggerDoc)

	paths, ok := swaggerDoc["paths"].(map[string]interface{})
	if !ok {
		t.Fatal("Поле 'paths' не найдено или имеет неверный тип")
	}

	// Проверяем наличие ожидаемых эндпоинтов
	expectedEndpoints := []string{"/health", "/go/{path}", "/python/{path}"}
	for _, ep := range expectedEndpoints {
		if _, exists := paths[ep]; !exists {
			t.Errorf("Эндпоинт %s не найден в Swagger документации", ep)
		}
	}
}

func TestSwaggerInfo(t *testing.T) {
	// Запускаем сервер в горутине
	go main()
	time.Sleep(2 * time.Second)

	resp, err := http.Get("http://localhost:8082/swagger/doc.json")
	if err != nil {
		t.Fatalf("Не удалось подключиться к doc.json: %v", err)
	}
	defer resp.Body.Close()

	var swaggerDoc map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&swaggerDoc)

	info, ok := swaggerDoc["info"].(map[string]interface{})
	if !ok {
		t.Fatal("Поле 'info' не найдено")
	}

	// Проверяем заголовок API
	if title, ok := info["title"]; ok {
		if title != "API Gateway" {
			t.Errorf("Заголовок API '%s', ожидался 'API Gateway'", title)
		}
	} else {
		t.Error("В документации отсутствует поле 'title'")
	}
}
