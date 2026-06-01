package scheduler

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"prayer-bot/internal/telegram"
)

func Start(
	bot *telegram.Bot,
	db *pgxpool.Pool,
) {
	log.Println("Scheduler started")

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkNotifications(bot, db)
		}
	}
}

func checkNotifications(
	bot *telegram.Bot,
	db *pgxpool.Pool,
) {
	log.Println("Checking prayer notifications...")
}
