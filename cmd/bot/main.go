package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"prayer-bot/internal/config"
	"prayer-bot/internal/database"
	"prayer-bot/internal/scheduler"
	"prayer-bot/internal/telegram"
)
func startHealthServer() {
 port := os.Getenv("PORT")
 if port == "" {
  port = "10000"
 }

 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Prayer Bot Running"))
 })

 go func() {
  log.Printf("Health server listening on :%s", port)
  if err := http.ListenAndServe(":"+port, nil); err != nil {
   log.Fatal(err)
  }
 }()
}
func main() {
	log.Println("Starting Prayer Bot...")
	startHealthServer()
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	log.Println("Database connected")

	// Create telegram bot
	bot, err := telegram.NewBot(cfg.BotToken, db)
	if err != nil {
		log.Fatalf("failed to create telegram bot: %v", err)
	}

	// Start scheduler
	go scheduler.Start(bot, db)

	// Start telegram polling
	go bot.Start()

	log.Println("Prayer Bot is running")

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	log.Println("Shutting down Prayer Bot...")
}
