import unittest
import requests
import subprocess
import time
import sys
import os

class TestAPIGateway(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        """Запускаем Go-сервис, FastAPI и Gateway перед тестами"""
        
        # Запускаем Go-сервис
        go_path = os.path.join(os.path.dirname(__file__), "..", "go_service")
        cls.go_server = subprocess.Popen(
            ["go", "run", "main.go"],
            cwd=go_path,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        
        # Запускаем FastAPI
        fastapi_path = os.path.join(os.path.dirname(__file__), "..", "python_service")
        cls.fastapi_server = subprocess.Popen(
            [sys.executable, "-m", "uvicorn", "main:app", "--port", "8000"],
            cwd=fastapi_path,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        
        # Запускаем Gateway
        cls.gateway = subprocess.Popen(
            ["go", "run", "gateway.go"],
            cwd=go_path,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        
        time.sleep(5)  # Ждём, пока всё запустится

    @classmethod
    def tearDownClass(cls):
        """Останавливаем все сервисы"""
        cls.gateway.terminate()
        cls.fastapi_server.terminate()
        cls.go_server.terminate()
        cls.gateway.wait()
        cls.fastapi_server.wait()
        cls.go_server.wait()

    def test_health_endpoint(self):
        """Проверяем health-эндпоинт Gateway"""
        response = requests.get("http://localhost:8082/health")
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.json()["status"], "ok")

    def test_go_proxy(self):
        """Проверяем прокси на Go-сервис"""
        response = requests.get("http://localhost:8082/go/")
        self.assertEqual(response.status_code, 200)
        self.assertIn("Hello from Go", response.text)

    def test_python_proxy(self):
        """Проверяем прокси на FastAPI"""
        response = requests.get("http://localhost:8082/python/call-go")
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertEqual(data["status"], "success")

if __name__ == "__main__":
    unittest.main()