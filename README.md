# Лабораторная работа №10: Веб-разработка — FastAPI (Python) vs Gin (Go)
**Автор:** Холхалова Алина Сергеевна  
**Группа:** 220032-11  
**Вариант:** 8
---
## Цель работы

Изучить и сравнить два подхода к созданию веб-сервисов на Python (FastAPI) и Go (Gin). Реализовать взаимодействие между сервисами на разных языках через REST и WebSocket, освоить инструменты документирования API (Swagger/OpenAPI).

---
## Выполненные задания
### Средняя сложность
| № | Задание | Статус |
|---|---------|--------|
| 2 | Добавить middleware для логирования в Go | ✅ |
| 4 | Создать FastAPI-сервис, который вызывает Go-сервис через HTTP | ✅ |
| 8 | Добавить Swagger-документацию для FastAPI и OpenAPI для Gin | ✅ |

### Повышенная сложность
| № | Задание | Статус |
|---|---------|--------|
| 2 | Создать API-шлюз на Go, который маршрутизирует запросы к разным микросервисам (Python, Go) | ✅ |
| 4 | Использовать WebSocket: реализовать чат на Go и подключиться к нему из Python | ✅ |
---
## Краткое описание проекта
Проект представляет собой систему из пяти взаимодействующих сервисов:

1. **Go-сервис (порт 8081)** — HTTP-сервер с middleware для логирования всех входящих запросов. Возвращает приветственное сообщение.

2. **FastAPI-сервис (порт 8000)** — Python-сервис, который при обращении к эндпоинту `/call-go` отправляет HTTP-запрос к Go-сервису и возвращает его ответ. Имеет встроенную Swagger-документацию.

3. **API Gateway (порт 8082)** — маршрутизатор на Go (Gin), который проксирует запросы:
   - `/go/*` → Go-сервис (порт 8081)
   - `/python/*` → FastAPI-сервис (порт 8000)
   - `/health` — проверка работоспособности шлюза

4. **WebSocket чат-сервер (порт 8083)** — сервер на Go с использованием библиотеки Gorilla WebSocket. Поддерживает множество подключений и рассылает сообщения всем участникам чата.

5. **Python-клиент для чата** — скрипт на Python с библиотекой `websockets`, который подключается к чат-серверу, отправляет и получает сообщения в реальном времени.
---
## Технологии
| Компонент | Технологии |
|-----------|------------|
| **Go-сервисы** | Go 1.22+, Gin, net/http, Gorilla WebSocket, httputil |
| **Python-сервисы** | Python 3.13, FastAPI, Uvicorn, Requests, Websockets |
| **Документация** | Swagger UI (FastAPI), Swaggo/OpenAPI (Gin) |
| **Тестирование** | unittest (Python), testing (Go) |
---
## Структура проекта
lab10/
├── go_service/ # Go-сервисы
│ ├── main.go # HTTP сервер с логированием (порт 8081)
│ ├── gateway.go # API Gateway (порт 8082)
│ ├── main_test.go # Go тесты для middleware
│ ├── gateway_test.go # Go тесты для Gateway
│ └── openapi_test.go # Go тесты для OpenAPI
├── chat_server/ # WebSocket чат на Go
│ ├── main.go # Чат сервер (порт 8083)
│ └── main_test.go # Go тесты для чата
├── python_service/ # FastAPI сервис
│ └── main.py # FastAPI приложение (порт 8000)
├── python_client/ # Python клиент для чата
│ └── chat_client.py # WebSocket клиент
├── tests/ # Python интеграционные тесты
│ ├── test_go_service.py # Тест Go-сервиса
│ ├── test_fastapi.py # Тест FastAPI
│ ├── test_gateway.py # Тест API Gateway
│ ├── test_websocket.py # Тест WebSocket чата
│ └── test_swagger.py # Тест документации
├── requirements.txt # Python зависимости
├── .gitignore # Игнорируемые файлы
├── README.md # Документация
└── PROMPT_LOG.md # История промптов и ошибок

---
## Запуск проекта
### Требования
- **Go** 1.22 или выше
- **Python** 3.10 или выше
- **Git**
## Установка зависимостей
### Python зависимости
pip install -r requirements.txt
### Go зависимости
cd go_service
go mod init go_service
go get github.com/gin-gonic/gin
cd ../chat_server
go mod init chat_server
go get github.com/gorilla/websocket
cd ..
## Запуск сервисов
**Важно:** каждый сервис запускается в отдельном терминале.
| Терминал | Команда | Порт | Назначение |
|----------|---------|------|------------|
| 1 | `cd go_service && go run main.go` | 8081 | Go-сервис с логированием |
| 2 | `cd python_service && uvicorn main:app --reload --port 8000` | 8000 | FastAPI сервис |
| 3 | `cd go_service && go run gateway.go` | 8082 | API Gateway |
| 4 | `cd chat_server && go run main.go` | 8083 | WebSocket чат |
| 5 | `cd python_client && python chat_client.py` | — | Клиент чата |
---
## Примеры запросов
### 1. Проверка Go-сервиса
curl http://localhost:8081/
Ответ: Hello from Go with logging!
### 2. FastAPI вызывает Go-сервис
curl http://localhost:8000/call-go
Ответ:
  "status": "success",
  "go_response": "Hello from Go with logging!"
### 3. API Gateway
### Проверка здоровья шлюза
curl http://localhost:8082/health
### Прокси на Go-сервис
curl http://localhost:8082/go/
### Прокси на FastAPI
curl http://localhost:8082/python/call-go
### 4. Документация API
Сервис	Адрес
FastAPI Swagger UI	http://localhost:8000/docs
FastAPI OpenAPI JSON	http://localhost:8000/openapi.json
Gin OpenAPI (если настроен)	http://localhost:8082/swagger/index.html
### 5. WebSocket чат
Запустите два экземпляра клиента в разных терминалах:

cd python_client
python chat_client.py
В первом клиенте введите сообщение — оно мгновенно появится во втором.

## Тестирование
### Python тесты
python -m unittest discover tests -v
Результат: 12 тестов, все OK

### Go тесты
cd go_service
go test -v

cd ../chat_server
go test -v
Результат: PASS

---
## Выводы
FastAPI позволяет быстро разрабатывать веб-сервисы благодаря встроенной валидации, автодокументации и простоте синтаксиса. Подходит для прототипирования и проектов, где скорость разработки важнее производительности.

Gin (Go) демонстрирует более высокую производительность и меньший размер бинарных файлов. Middleware легко реализуются через функции-обёртки. Go предпочтителен для высоконагруженных микросервисов.

API Gateway на Go успешно маршрутизирует запросы между сервисами, добавляя единую точку входа и централизованное логирование.

WebSocket на Go работает стабильно и выдерживает множество одновременных подключений. Python-клиент легко интегрируется через библиотеку websockets.

Все задания выполнены в полном объёме, код покрыт тестами (Python + Go), документация оформлена.

---
## Используемые технологии
Go 1.22 + Gin + Gorilla WebSocket + net/http/httputil

Python 3.13 + FastAPI + Uvicorn + Requests + websockets

Swagger UI / OpenAPI / Swaggo

unittest / testing

## Ссылка на репозиторий
https://github.com/holhalovaa/LR_10