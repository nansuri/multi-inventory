package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"multi-inventory/internal/application"
	httpHandler "multi-inventory/internal/infrastructure/http"
	"multi-inventory/internal/infrastructure/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to Database
	db, err := postgres.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Auto-migrate tables on startup (MVP-friendly)
	if err := db.AutoMigrate(context.Background()); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	fmt.Println("AutoMigrate executed successfully.")

	// Initialize Layers
	userRepo := postgres.NewUserRepository(db)
	itemRepo := postgres.NewItemRepository(db)
	orderRepo := postgres.NewOrderRepository(db)

	authService := application.NewAuthService(userRepo)
	inventoryService := application.NewInventoryService(itemRepo)
	salesService := application.NewSalesService(orderRepo, itemRepo)

	authHandler := httpHandler.NewAuthHandler(authService)
	inventoryHandler := httpHandler.NewInventoryHandler(inventoryService)
	salesHandler := httpHandler.NewSalesHandler(salesService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Multi Inventory API is running!"))
	})

	r.Mount("/api/auth", authHandler.Routes())
	r.Mount("/api/inventory", inventoryHandler.Routes())
	r.Mount("/api/sales", salesHandler.Routes())

	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
