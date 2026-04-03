import unittest
import requests
import subprocess
import time
import sys
import os

class TestGoService(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        """Запускаем Go-сервис перед тестами"""
        go_path = os.path.join(os.path.dirname(__file__), "..", "go_service")
        cls.server = subprocess.Popen(
            ["go", "run", "main.go"],
            cwd=go_path,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        time.sleep(3)  # Ждём, пока сервер запустится

    @classmethod
    def tearDownClass(cls):
        """Останавливаем сервер после тестов"""
        cls.server.terminate()
        cls.server.wait()

    def test_hello_endpoint(self):
        """Проверяем, что эндпоинт / возвращает Hello"""
        response = requests.get("http://localhost:8081/")
        self.assertEqual(response.status_code, 200)
        self.assertIn("Hello from Go", response.text)

    def test_server_running(self):
        """Проверяем, что сервер вообще отвечает"""
        response = requests.get("http://localhost:8081/")
        self.assertEqual(response.status_code, 200)

if __name__ == "__main__":
    unittest.main()