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
	"WB-TechSchool-L0/pkg/db"
)

func main() {
	// --- подключение к БД ---
	pgDb, err := db.ConnectToDb()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer pgDb.Close()

	// --- инициализация repo ---
	orderRepo := repo.NewPgOrderRepo(pgDb)

	// --- чтение тестового заказа из файла model.json ---
	data, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatal("Ошибка чтения model.json:", err)
	}

	var order domain.Order
	if err := json.Unmarshal(data, &order); err != nil {
		log.Fatal("Ошибка парсинга JSON:", err)
	}

	// вставка даты, если её нет
	if order.DateCreated.IsZero() {
		order.DateCreated = time.Now()
	}

	// --- сохранение заказа ---
	ctx := context.Background()
	if err := orderRepo.Save(ctx, &order); err != nil {
		log.Fatal("Ошибка сохранения заказа:", err)
	}
	fmt.Println("✅ Заказ сохранён:", order.OrderUid)

	// --- достаём обратно ---
	got, err := orderRepo.GetById(ctx, order.OrderUid)
	if err != nil {
		log.Fatal("Ошибка получения заказа:", err)
	}

	fmt.Println("✅ Заказ получен из БД:")
	fmt.Printf("%+v\n", got)
}
