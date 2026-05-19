# 🛠️ Go Developer REST API Service

A production-ready RESTful CRUD API service written in Go. This repository demonstrates modern software engineering practices, clean code architecture, and idiomatic Go development patterns.

The project is designed to showcase enterprise-level backend development fundamentals, focusing on reliability, data validation, and clean database integration.

## 🚀 Key Architectural Features

* **Clean Architecture / Layered Pattern:** Clear separation of concerns between HTTP handlers, business logic (use cases/services), and data access layers (repositories).
* **RESTful Standards:** Full implementation of CRUD operations with proper HTTP status codes, structured JSON payloads, and error handling.
* **Database Integration:** Secure database connectivity featuring connection pooling and structured migrations.
* **Input Validation:** Robust request payload verification to enforce data integrity before hitting the business layer.
* **Environment Configuration:** Twelve-Factor App compliance using environment variables for safe configuration management.

## 🛠️ Tech Stack

* **Language:** Go (Golang)
* **API / Routing:** Idiomatic HTTP routing
* **Database:** PostgreSQL / MySQL (Structured relational storage)
* **Dependency Management:** Go Modules (`go.mod`)

## 🏁 Getting Started

### Prerequisites

* Go 1.18+
* Database instance configured (PostgreSQL/MySQL)

### Configuration

Create a `.env` file in the root directory or set up your environment variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_db_name
SERVER_PORT=8080
```

### Installation & Run

1. Clone the repository:
```bash
git clone https://github.com
cd go-developer-test-task
```

2. Download dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

## 📝 License

This project is open-source and available under the [MIT License](LICENSE).

---

# 📦 Subscription Service

REST API сервис для управления подписками пользователей и расчёта их суммарной стоимости за выбранный период.

---

## 🏗 Архитектура

Проект реализован с использованием **Clean Architecture (слоистая архитектура)**.

### Структура слоёв:

cmd/app → точка входа  
internal/domain → бизнес-сущности и правила  
internal/dto → входные DTO  
internal/handler → HTTP слой  
internal/service → бизнес-логика  
internal/repository → интерфейсы репозиториев  
internal/repository/postgres → реализация работы с БД  
migrations → SQL миграции


### Принципы:

- Разделение ответственности (SRP)
- Зависимость направлена внутрь (Dependency Inversion)
- Бизнес-логика не зависит от инфраструктуры
- Repository Pattern
- DTO отделены от доменных моделей

---

## 🚀 Функциональность

- ✅ Создание подписки
- ✅ Получение подписки по ID
- ✅ Обновление подписки
- ✅ Удаление подписки
- ✅ Получение списка подписок
    - по пользователю
    - по сервису
- ✅ Подсчёт суммарной стоимости подписок за период с фильтрами:
    - по пользователю
    - по сервису

---

## 🛠 Используемые технологии

- Go 1.24
- PostgreSQL 16
- pgx/v5
- chi router
- Docker + docker-compose
- golang-migrate

---

## 📅 Формат дат

В API используется формат:

MM-YYYY


Пример:

01-2025  
12-2024  
02-2026


---

## ⚙️ Запуск проекта

### 1️⃣ Склонировать репозиторий

```bash
git clone <repo_url>
cd <project>
```

### 2️⃣ Создать .env
```bash
cp .env.example .env
```

Пример .env.example:  
APP_PORT=8081

DB_HOST=db  
DB_PORT=5432  
DB_USER=postgres  
DB_PASSWORD=secret  
DB_NAME=subscriptions  
DB_SSLMODE=disable

### 3️⃣ Запустить
```bash
docker compose up --build
```
Сервис будет доступен по адресу:  
http://localhost:8081  
Адрес swagger:  
http://localhost:8081/swagger/index.html#/

## 📌 API Endpoints
Создание подписки
```bash
POST /subscriptions
```
Получение подписки
```bash
GET /subscriptions/{id}
```
Обновление подписки
```bash
PATCH /subscriptions/{id}
```
Удаление подписки
```bash
DELETE /subscriptions/{id}
```
Список подписок
```bash
GET /subscriptions/list?user_id=&service_name=
```
Подсчёт общей стоимости
```bash
GET /subscriptions/total?from=01-2025&to=12-2025
GET /subscriptions/total?user_id=UUID&service_name=ServiceName&from=MM-YYYY&to=MM-YYYY
```

## 🗄 База данных
Используется PostgreSQL.

Особенности:
- UUID в качестве primary key
- CHECK constraint для price >= 0

Индексы:
- по user_id
- по service_name
- по датам подписки

Миграции выполняются автоматически при запуске контейнера.

## 📐 Принятые решения
- UUID вместо автоинкремента — унификация типов идентификаторов.
- Разделение DTO и domain моделей.
- Валидация доменной сущности в service слое.
- Использование контекста во всех слоях.
- Middleware для логирования HTTP-запросов.

## 🧪 Возможные улучшения
- Централизованная обработка ошибок
- Структурированное логирование (zap / slog)
- Добавление unit-тестов
- Использование интерфейсов для service слоя
- Pagination для списка подписок
- Метрики (Prometheus)

## 📄 Лицензия
Test task project.
