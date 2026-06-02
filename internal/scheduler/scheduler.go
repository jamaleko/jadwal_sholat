package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"prayer-bot/internal/service"
	"prayer-bot/internal/telegram"
	tgbot "github.com/go-telegram/bot"
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
	ctx := context.Background()

	log.Println("Checking prayer notifications...")

	userService := service.NewUserService(db)
	notifService := service.NewNotificationService(db)
	prayerService := service.NewPrayerService()

	users, err := userService.GetActiveUsers(ctx)
	if err != nil {
		log.Printf("failed get active users: %v", err)
		return
	}

	now := time.Now()
	today := now
	//testSchedule := true
	

	for _, user := range users {
		schedule, err := prayerService.GetTodaySchedule(
			user.Latitude,
			user.Longitude,
		)
		if err != nil {
			log.Printf("failed get prayer schedule for chatID=%d: %v", user.ChatID, err)
			continue
		}
		/*if testSchedule {
        // override semua jadwal menjadi 1 menit dari sekarang
        schedule.Dhuhr = now
    }*/
		checkAndSend(bot, notifService, user.ChatID, "Subuh", schedule.Fajr, today)
		checkAndSend(bot, notifService, user.ChatID, "Dzuhur", schedule.Dhuhr, today)
		checkAndSend(bot, notifService, user.ChatID, "Ashar", schedule.Asr, today)
		checkAndSend(bot, notifService, user.ChatID, "Maghrib", schedule.Maghrib, today)
		checkAndSend(bot, notifService, user.ChatID, "Isya", schedule.Isha, today)
	}
}

func checkAndSend(
 bot *telegram.Bot,
 notifService *service.NotificationService,
 chatID int64,
 prayerName string,
 prayerTime time.Time,
 today time.Time,
) {
 ctx := context.Background()

 loc, err := time.LoadLocation("Asia/Jakarta")
 if err != nil {
  log.Printf("failed load timezone: %v", err)
  return
 }

 now := time.Now().In(loc)
 prayerTime = prayerTime.In(loc)

 /*log.Printf(
  "NOW=%s TARGET=%s PRAYER=%s",
  now.Format("2006-01-02 15:04:05"),
  prayerTime.Format("2006-01-02 15:04:05"),
  prayerName,
 )*/

 // cek jam dan menit saja
 if now.Hour() != prayerTime.Hour() ||
  now.Minute() != prayerTime.Minute() {
  return
 }

 sent, err := notifService.AlreadySent(
  ctx,
  chatID,
  prayerName,
  now,
 )

 if err != nil {
  log.Printf(
   "failed check already sent: %v",
   err,
  )
  return
 }

 if sent {
  return
 }

 msg := fmt.Sprintf(
  "🔔 Waktu %s telah tiba! ⏰",
  prayerName,
 )

 _, err = bot.Client().SendMessage(
  ctx,
  &tgbot.SendMessageParams{
   ChatID: chatID,
   Text:   msg,
  },
 )

 if err != nil {
  log.Printf(
   "failed send notification to chatID=%d: %v",
   chatID,
   err,
  )
  return
 }

 err = notifService.MarkAsSent(
  ctx,
  chatID,
  prayerName,
  now,
 )

 if err != nil {
  log.Printf(
   "failed mark notification sent: %v",
   err,
  )
 }
}
