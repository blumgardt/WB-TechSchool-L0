# WB-TechSchool-L0

Демонстрационный сервис заказов на Go с использованием **PostgreSQL**, **Redis** и **Kafka**.  
Сервис получает заказы из Kafka, сохраняет их в Postgres, кэширует в Redis и отдаёт через HTTP API и простую веб-страницу.

---

## 📂 Структура проекта
WB-TechSchool-L0/
├── cmd/                # Точка входа (main.go)
├── internal/
│   ├── domain/         # Бизнес-модели (Order, Delivery, Payment, Item)
│   ├── repo/           # Репозиторий (Postgres реализация)
│   ├── service/
│   │   ├── redis/      # Работа с Redis
│   │   └── kafka/      # Kafka consumer
│   └── http/           # HTTP-роутер и хендлеры
├── pkg/
│   └── db/             # Подключение к БД, миграции
├── web/static/         # HTML/JS интерфейс
├── go.mod
├── go.sum
├── README.md
└── .env.example

---

## 🛠 Используемые технологии

- Go 1.24+
- PostgreSQL 15+
- Redis 7+
- Kafka (Confluent)
- sqlx
- go-redis
- segmentio/kafka-go

---
