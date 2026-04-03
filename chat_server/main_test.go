package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketUpgrade(t *testing.T) {
	// Создаём тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(handleConnections))
	defer server.Close()

	// Преобразуем http:// в ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Подключаемся к WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Не удалось подключиться к WebSocket: %v", err)
	}
	defer ws.Close()

	// Проверяем, что соединение установлено
	if ws == nil {
		t.Error("WebSocket соединение не установлено")
	}
}

func TestBroadcastMessage(t *testing.T) {
	// Создаём тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(handleConnections))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Подключаем первого клиента
	ws1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Первый клиент не подключился: %v", err)
	}
	defer ws1.Close()

	// Подключаем второго клиента
	ws2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Второй клиент не подключился: %v", err)
	}
	defer ws2.Close()

	// Отправляем сообщение от первого клиента
	testMsg := "Hello from test"
	err = ws1.WriteMessage(websocket.TextMessage, []byte(testMsg))
	if err != nil {
		t.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	// Ждём и читаем сообщение у второго клиента
	done := make(chan bool)
	go func() {
		_, msg, err := ws2.ReadMessage()
		if err != nil {
			t.Errorf("Ошибка чтения сообщения: %v", err)
			done <- false
			return
		}
		if string(msg) != testMsg {
			t.Errorf("Получено неправильное сообщение: получили %s, ожидали %s", msg, testMsg)
			done <- false
			return
		}
		done <- true
	}()

	// Таймаут на ожидание
	select {
	case result := <-done:
		if !result {
			t.Error("Тест не прошёл")
		}
	case <-time.After(5 * time.Second):
		t.Error("Таймаут: сообщение не получено за 5 секунд")
	}
}

func TestMultipleClients(t *testing.T) {
	// Создаём тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(handleConnections))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Подключаем 3 клиентов
	clients := make([]*websocket.Conn, 3)
	for i := 0; i < 3; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("Клиент %d не подключился: %v", i, err)
		}
		defer ws.Close()
		clients[i] = ws
	}

	// Отправляем сообщение от первого клиента
	testMsg := "Broadcast test"
	err := clients[0].WriteMessage(websocket.TextMessage, []byte(testMsg))
	if err != nil {
		t.Fatalf("Ошибка отправки: %v", err)
	}

	// Проверяем, что все остальные получили сообщение
	for i := 1; i < 3; i++ {
		_, msg, err := clients[i].ReadMessage()
		if err != nil {
			t.Errorf("Клиент %d не получил сообщение: %v", i, err)
			continue
		}
		if string(msg) != testMsg {
			t.Errorf("Клиент %d получил: %s, ожидал: %s", i, msg, testMsg)
		}
	}
}
