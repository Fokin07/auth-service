package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AlexFox86/auth-service/internal/delivery"
	"github.com/AlexFox86/auth-service/internal/repository/postgres"
	"github.com/AlexFox86/auth-service/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
}

func main() {
	db, err := sqlx.Connect("postgres", connectionString())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := postgres.NewPgRepository(db)
	service := service.New(repo, os.Getenv("SECRET"), time.Hour)
	handler := delivery.NewHandler(service)

	http.HandleFunc("POST /register", handler.Register)
	http.HandleFunc("POST /login", handler.Login)
	http.HandleFunc("GET /validate", handler.Validate)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
