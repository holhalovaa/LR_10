package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	// Создаём запрос
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаём ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Создаём handler с middleware
	handler := loggingMiddleware(helloHandler)

	// Вызываем handler
	handler.ServeHTTP(rr, req)

	// Проверяем статус
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler вернул неправильный статус: получили %v, ожидали %v",
			status, http.StatusOK)
	}

	// Проверяем тело ответа
	expected := "Hello from Go with logging!"
	if rr.Body.String() != expected {
		t.Errorf("handler вернул неправильное тело: получили %v, ожидали %v",
			rr.Body.String(), expected)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	// Просто проверяем, что middleware не падает
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := loggingMiddleware(helloHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("middleware не сработал корректно")
	}
}
