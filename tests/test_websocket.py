import unittest
import subprocess
import time
import sys
import os
import threading
import websocket
import json

class TestWebSocketChat(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        """Запускаем WebSocket сервер перед тестами"""
        chat_path = os.path.join(os.path.dirname(__file__), "..", "chat_server")
        cls.server = subprocess.Popen(
            ["go", "run", "main.go"],
            cwd=chat_path,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        time.sleep(3)
        
        cls.messages_received = []
        
    @classmethod
    def tearDownClass(cls):
        """Останавливаем сервер"""
        cls.server.terminate()
        cls.server.wait()

    def on_message(self, ws, message):
        """Колбэк для получения сообщений"""
        self.messages_received.append(message)

    def test_websocket_connection(self):
        """Проверяем, что можно подключиться к WebSocket"""
        ws = websocket.WebSocket()
        ws.connect("ws://localhost:8083/ws")
        self.assertTrue(ws.connected)
        ws.close()

    def test_send_and_receive(self):
        """Проверяем, что сообщение отправляется и получается"""
        received = []
        
        def on_message(ws, message):
            received.append(message)
        
        # Создаём два клиента
        ws1 = websocket.WebSocketApp("ws://localhost:8083/ws",
                                     on_message=on_message)
        ws2 = websocket.WebSocketApp("ws://localhost:8083/ws",
                                     on_message=on_message)
        
        # Запускаем клиентов в потоках
        thread1 = threading.Thread(target=ws1.run_forever)
        thread2 = threading.Thread(target=ws2.run_forever)
        thread1.daemon = True
        thread2.daemon = True
        thread1.start()
        thread2.start()
        
        time.sleep(1)
        
        # Отправляем сообщение через второго клиента
        ws2.send("Hello from test")
        time.sleep(1)
        
        # Проверяем, что первый клиент получил сообщение
        self.assertIn("Hello from test", str(received))
        
        ws1.close()
        ws2.close()

if __name__ == "__main__":
    unittest.main()