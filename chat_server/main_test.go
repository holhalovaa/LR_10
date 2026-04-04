package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketConnection(t *testing.T) {
	port := 8084
	srv := RunServer(port)
	defer srv.Close()

	time.Sleep(1 * time.Second)

	ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%d/ws", port), nil)
	if err != nil {
		t.Fatalf("Не удалось подключиться: %v", err)
	}
	defer ws.Close()

	if ws == nil {
		t.Error("Соединение не установлено")
	}
}

func TestSendAndReceiveMessage(t *testing.T) {
	port := 8085
	srv := RunServer(port)
	defer srv.Close()

	time.Sleep(1 * time.Second)

	ws1, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%d/ws", port), nil)
	if err != nil {
		t.Fatalf("Первый клиент не подключился: %v", err)
	}
	defer ws1.Close()

	ws2, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%d/ws", port), nil)
	if err != nil {
		t.Fatalf("Второй клиент не подключился: %v", err)
	}
	defer ws2.Close()

	time.Sleep(500 * time.Millisecond)

	testMsg := "Hello, chat!"
	err = ws1.WriteMessage(websocket.TextMessage, []byte(testMsg))
	if err != nil {
		t.Fatalf("Ошибка отправки: %v", err)
	}

	_, msgBytes, err := ws2.ReadMessage()
	if err != nil {
		t.Fatalf("Ошибка чтения: %v", err)
	}

	receivedMsg := string(msgBytes)
	if receivedMsg != testMsg {
		t.Errorf("Получено: %s, ожидалось: %s", receivedMsg, testMsg)
	}
}

func TestMultipleClientsBroadcast(t *testing.T) {
	port := 8086
	srv := RunServer(port)
	defer srv.Close()

	time.Sleep(1 * time.Second)

	clients := make([]*websocket.Conn, 3)
	for i := 0; i < 3; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%d/ws", port), nil)
		if err != nil {
			t.Fatalf("Клиент %d не подключился: %v", i, err)
		}
		defer ws.Close()
		clients[i] = ws
	}

	time.Sleep(500 * time.Millisecond)

	testMsg := "Broadcast message"
	err := clients[0].WriteMessage(websocket.TextMessage, []byte(testMsg))
	if err != nil {
		t.Fatalf("Ошибка отправки: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	for i := 1; i < 3; i++ {
		_, msgBytes, err := clients[i].ReadMessage()
		if err != nil {
			t.Errorf("Клиент %d ошибка чтения: %v", i, err)
			continue
		}
		if string(msgBytes) != testMsg {
			t.Errorf("Клиент %d получил: %s, ожидал: %s", i, msgBytes, testMsg)
		}
	}
}
