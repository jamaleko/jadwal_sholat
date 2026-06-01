package telegram

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) registerHandlers() {
	b.client.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"/start",
		tgbot.MatchTypeExact,
		b.handleStart,
	)

	b.client.RegisterHandler(
		tgbot.HandlerTypeMessageLocation,
		"",
		tgbot.MatchTypeExact,
		b.handleLocation,
	)
}

func (b *Bot) handleStart(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
) {
	chatID := update.Message.Chat.ID

	msg := `
Assalamu'alaikum 👋

Selamat datang di Prayer Bot.

Untuk mulai menerima notifikasi jadwal sholat, silakan bagikan lokasi Anda menggunakan tombol di bawah.
`

	_, err := bot.SendMessage(
		ctx,
		&tgbot.SendMessageParams{
			ChatID: chatID,
			Text:   msg,
			ReplyMarkup: MainKeyboard(),
		},
	)

	if err != nil {
		log.Printf(
			"failed send start message: %v",
			err,
		)
	}
}

func (b *Bot) handleLocation(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
) {
	location := update.Message.Location

	log.Printf(
		"location received: lat=%f lon=%f",
		location.Latitude,
		location.Longitude,
	)

	_, err := bot.SendMessage(
		ctx,
		&tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: "📍 Lokasi berhasil diterima.\n\n" +
				"Notifikasi jadwal sholat akan segera diaktifkan.",
		},
	)

	if err != nil {
		log.Printf(
			"failed send location confirmation: %v",
			err,
		)
	}
}
