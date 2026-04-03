package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка подключения: %v", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	log.Printf("✅ Клиент подключен. Всего клиентов: %d", len(clients))

	for {
		// Читаем сообщение как ТЕКСТ (не JSON)
		_, msgBytes, err := ws.ReadMessage()
		if err != nil {
			log.Printf("❌ Клиент отключился: %v", err)
			delete(clients, ws)
			break
		}
		msg := string(msgBytes)
		log.Printf("📩 Получено сообщение: '%s'", msg)
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		log.Printf("📢 Рассылаю сообщение '%s' всем %d клиентам", msg, len(clients))
		for client := range clients {
			// Отправляем как ТЕКСТ
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("⚠️ Ошибка отправки: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	fmt.Println("🚀 Чат сервер запущен на порту 8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
