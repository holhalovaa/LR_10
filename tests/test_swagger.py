import unittest
import requests
import subprocess
import time
import sys
import os

class TestSwaggerDocumentation(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        """Запускаем FastAPI сервис перед тестами"""
        fastapi_path = os.path.join(os.path.dirname(__file__), "..", "python_service")
        cls.server = subprocess.Popen(
            [sys.executable, "-m", "uvicorn", "main:app", "--port", "8000"],
            cwd=fastapi_path,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        time.sleep(3)

    @classmethod
    def tearDownClass(cls):
        """Останавливаем сервер после тестов"""
        cls.server.terminate()
        cls.server.wait()

    def test_swagger_ui_page(self):
        """Проверяем, что Swagger UI доступен"""
        response = requests.get("http://localhost:8000/docs")
        self.assertEqual(response.status_code, 200)
        self.assertIn("text/html", response.headers.get("Content-Type", ""))

    def test_openapi_json(self):
        """Проверяем, что OpenAPI JSON схема доступна"""
        response = requests.get("http://localhost:8000/openapi.json")
        self.assertEqual(response.status_code, 200)
        
        data = response.json()
        # Проверяем, что в схеме есть информация о нашем API
        self.assertIn("info", data)
        self.assertIn("title", data["info"])
        self.assertEqual(data["info"]["title"], "Python Service")
        
        # Проверяем, что эндпоинт /call-go описан в схеме
        self.assertIn("/call-go", data["paths"])

    def test_swagger_title(self):
        """Проверяем, что заголовок Swagger соответствует заданию"""
        response = requests.get("http://localhost:8000/docs")
        self.assertEqual(response.status_code, 200)
        # Заголовок должен содержать название сервиса
        self.assertIn("Python Service", response.text)

if __name__ == "__main__":
    unittest.main()