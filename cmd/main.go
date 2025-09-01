package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"WB-TechSchool-L0/internal/domain"
	"WB-TechSchool-L0/internal/repo"
	"WB-TechSchool-L0/internal/service"
	"WB-TechSchool-L0/internal/service/redis"
	"WB-TechSchool-L0/pkg/db"
)

func main() {
	ctx := context.Background()

	// --- подключение к БД ---
	pgDb, err := db.ConnectToDb()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer pgDb.Close()

	// --- инициализация repo ---
	orderRepo := repo.NewPgOrderRepo(pgDb)

	// --- инициализация Redis ---
	redisClient := redis.NewRedisClient("localhost:6379", 5*time.Minute, orderRepo)
	if err := redisClient.RestoreCache(ctx); err != nil {
		log.Fatal("Ошибка восстановления кэша:", err)
	}

	// --- инициализация сервиса ---
	orderService := service.NewOrderService(orderRepo, redisClient)

	// --- читаем тестовый заказ из файла model.json ---
	data, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatal("Ошибка чтения model.json:", err)
	}

	var order domain.Order
	if err := json.Unmarshal(data, &order); err != nil {
		log.Fatal("Ошибка парсинга JSON:", err)
	}

	if order.DateCreated.IsZero() {
		order.DateCreated = time.Now()
	}

	// --- сохраняем заказ в БД ---
	if err := orderRepo.Save(ctx, &order); err != nil {
		log.Fatal("Ошибка сохранения заказа:", err)
	}
	fmt.Println("✅ Заказ сохранён в БД:", order.OrderUid)

	// --- получаем заказ через сервис (Redis + Repo) ---
	got, err := orderService.GetOrderById(order.OrderUid, ctx)
	if err != nil {
		log.Fatal("Ошибка получения заказа:", err)
	}
	fmt.Println("✅ Первый запрос (из БД, закэширован):", got.OrderUid)

	// --- второй запрос (должен вытащить уже из Redis) ---
	gotCached, err := orderService.GetOrderById(order.OrderUid, ctx)
	if err != nil {
		log.Fatal("Ошибка получения заказа:", err)
	}
	fmt.Println("✅ Второй запрос (из Redis):", gotCached.OrderUid)
}
