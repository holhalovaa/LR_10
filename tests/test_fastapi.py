import unittest
import requests
import subprocess
import time
import sys
import os

class TestFastAPIService(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        """Запускаем Go-сервис и FastAPI сервис перед тестами"""
        
        # Запускаем Go-сервис (нужен для FastAPI)
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
        
        time.sleep(5)  # Ждём, пока оба сервера запустятся

    @classmethod
    def tearDownClass(cls):
        """Останавливаем все серверы после тестов"""
        cls.fastapi_server.terminate()
        cls.go_server.terminate()
        cls.fastapi_server.wait()
        cls.go_server.wait()

    def test_call_go_endpoint(self):
        """Проверяем, что /call-go возвращает ответ"""
        response = requests.get("http://localhost:8000/call-go")
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertEqual(data["status"], "success")
        self.assertIn("Hello from Go", data["go_response"])

    def test_docs_endpoint(self):
        """Проверяем, что Swagger документация доступна"""
        response = requests.get("http://localhost:8000/docs")
        self.assertEqual(response.status_code, 200)

if __name__ == "__main__":
    unittest.main()