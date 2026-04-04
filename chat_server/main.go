package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.RWMutex // защита от race condition
	broadcast = make(chan string)
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка подключения: %v", err)
		return
	}
	defer ws.Close()

	// Блокируем мьютекс на запись при добавлении клиента
	clientsMu.Lock()
	clients[ws] = true
	clientsMu.Unlock()
	log.Printf("✅ Клиент подключен. Всего клиентов: %d", len(clients))

	for {
		_, msgBytes, err := ws.ReadMessage()
		if err != nil {
			log.Printf("❌ Клиент отключился: %v", err)
			// Блокируем мьютекс на запись при удалении клиента
			clientsMu.Lock()
			delete(clients, ws)
			clientsMu.Unlock()
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

		// Блокируем мьютекс на чтение при итерации по карте
		clientsMu.RLock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("⚠️ Ошибка отправки: %v", err)
				client.Close()
				// При удалении внутри итерации нужна запись, поэтому временно переключаемся
				clientsMu.RUnlock() // снимаем блокировку чтения
				clientsMu.Lock()    // берём блокировку записи
				delete(clients, client)
				clientsMu.Unlock() // отпускаем запись
				clientsMu.RLock()  // снова берём чтение для продолжения итерации
			}
		}
		clientsMu.RUnlock()
	}
}

func RunServer(port int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleConnections)

	go handleMessages()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		log.Printf("🚀 Чат сервер запущен на порту %d", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка сервера: %v", err)
		}
	}()

	return srv
}

func main() {
	RunServer(8083)
	select {} // Бесконечное ожидание
}
