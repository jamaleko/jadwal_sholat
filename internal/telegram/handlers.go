package telegram

import (
	"context"
	"fmt"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"prayer-bot/internal/service"
)

func (b *Bot) registerHandlers() {

	b.client.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"/start",
		tgbot.MatchTypeExact,
		b.handleStart,
	)

	b.client.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"📖 Jadwal Hari Ini",
		tgbot.MatchTypeExact,
		b.handleTodaySchedule,
	)
	
	b.client.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"🔕 Stop Notifikasi",
		tgbot.MatchTypeExact,
		b.handleDisableNotification,
	)
	
	b.client.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"🔔 Aktifkan Notifikasi",
		tgbot.MatchTypeExact,
		b.handleEnableNotification,
	)

	b.client.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {

			if update.Message == nil {
				return false
			}

			return update.Message.Location != nil
		},
		b.handleLocation,
	)
}
func (b *Bot) handleTodaySchedule(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
) {

	chatID := update.Message.Chat.ID

	userService := service.NewUserService(b.db)

	user, err := userService.GetByChatID(
		ctx,
		chatID,
	)

	if err != nil {

		_, _ = bot.SendMessage(
			ctx,
			&tgbot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ Silakan bagikan lokasi terlebih dahulu.",
			},
		)

		return
	}

	prayerService := service.NewPrayerService()

	schedule, err := prayerService.GetTodaySchedule(
		user.Latitude,
		user.Longitude,
	)

	if err != nil {

		_, _ = bot.SendMessage(
			ctx,
			&tgbot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ Gagal mengambil jadwal sholat.",
			},
		)

		return
	}

	message := fmt.Sprintf(
		`📖 Jadwal Sholat Hari Ini

🌅 Subuh   : %s
☀️ Dzuhur  : %s
🌤 Ashar   : %s
🌇 Maghrib : %s
🌙 Isya    : %s`,
		schedule.Fajr.Format("15:04"),
		schedule.Dhuhr.Format("15:04"),
		schedule.Asr.Format("15:04"),
		schedule.Maghrib.Format("15:04"),
		schedule.Isha.Format("15:04"),
	)

	_, _ = bot.SendMessage(
		ctx,
		&tgbot.SendMessageParams{
			ChatID: chatID,
			Text:   message,
		},
	)
}
func (b *Bot) handleDisableNotification(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
) {

	chatID := update.Message.Chat.ID

	userService := service.NewUserService(b.db)

	err := userService.Disable(
		ctx,
		chatID,
	)

	if err != nil {

		_, _ = bot.SendMessage(
			ctx,
			&tgbot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ Gagal menonaktifkan notifikasi.",
			},
		)

		return
	}

	_, _ = bot.SendMessage(
		ctx,
		&tgbot.SendMessageParams{
			ChatID: chatID,
			Text:   "🔕 Notifikasi berhasil dinonaktifkan.",
		},
	)
}
func (b *Bot) handleEnableNotification(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
) {

	chatID := update.Message.Chat.ID

	userService := service.NewUserService(b.db)

	err := userService.Enable(
		ctx,
		chatID,
	)

	if err != nil {

		_, _ = bot.SendMessage(
			ctx,
			&tgbot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ Gagal mengaktifkan notifikasi.",
			},
		)

		return
	}

	_, _ = bot.SendMessage(
		ctx,
		&tgbot.SendMessageParams{
			ChatID: chatID,
			Text:   "🔔 Notifikasi berhasil diaktifkan.",
		},
	)
}
func (b *Bot) handleStart(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
) {
	chatID := update.Message.Chat.ID

	msg := `Assalamu'alaikum 👋

Selamat datang di Jadwal Sholat by Tekuna.my.id

Silakan bagikan lokasi Anda untuk menerima notifikasi jadwal sholat.`

	_, err := bot.SendMessage(
		ctx,
		&tgbot.SendMessageParams{
			ChatID:     chatID,
			Text:       msg,
			ReplyMarkup: MainKeyboard(),
		},
	)

	if err != nil {
		log.Printf("failed send start message: %v", err)
	}
}

func (b *Bot) handleLocation(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
) {

	location := update.Message.Location
	chatID := update.Message.Chat.ID

	log.Printf(
		"location received: lat=%f lon=%f",
		location.Latitude,
		location.Longitude,
	)

	// Save user
	userService := service.NewUserService(
		b.db,
	)

	err := userService.SaveLocation(
		ctx,
		chatID,
		location.Latitude,
		location.Longitude,
	)

	if err != nil {

		log.Printf(
			"failed save user location: %v",
			err,
		)

		_, _ = bot.SendMessage(
			ctx,
			&tgbot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ Gagal menyimpan lokasi.",
			},
		)

		return
	}

	// Calculate prayer schedule
	prayerService := service.NewPrayerService()

	schedule, err := prayerService.GetTodaySchedule(
		location.Latitude,
		location.Longitude,
	)

	if err != nil {

		log.Printf(
			"failed calculate prayer schedule: %v",
			err,
		)

		_, _ = bot.SendMessage(
			ctx,
			&tgbot.SendMessageParams{
				ChatID: chatID,
				Text:   "❌ Gagal menghitung jadwal sholat.",
			},
		)

		return
	}

	message := fmt.Sprintf(
		`✅ Lokasi berhasil disimpan

Jadwal Sholat Hari Ini

🌅 Subuh   : %s
☀️ Dzuhur  : %s
🌤 Ashar   : %s
🌇 Maghrib : %s
🌙 Isya    : %s

🔔 Notifikasi telah diaktifkan.`,
		schedule.Fajr.Format("15:04"),
		schedule.Dhuhr.Format("15:04"),
		schedule.Asr.Format("15:04"),
		schedule.Maghrib.Format("15:04"),
		schedule.Isha.Format("15:04"),
	)

	_, err = bot.SendMessage(
		ctx,
		&tgbot.SendMessageParams{
			ChatID:     chatID,
			Text:       message,
			ReplyMarkup: MainKeyboard(),
		},
	)

	if err != nil {
		log.Printf(
			"failed send prayer schedule: %v",
			err,
		)
	}
}
