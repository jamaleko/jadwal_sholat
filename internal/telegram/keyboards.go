package telegram

import "github.com/go-telegram/bot/models"

func MainKeyboard() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]models.KeyboardButton{
			{
				{
					Text:            "📍 Bagikan Lokasi",
					RequestLocation: true,
				},
			},
			{
				{
					Text: "📖 Jadwal Hari Ini",
				},
				{
					Text: "🔕 Stop Notifikasi",
				},
			},
			{
				{
					Text: "🔔 Aktifkan Notifikasi",
				},
			},
		},
	}
}

func RemoveKeyboard() *models.ReplyKeyboardRemove {
	return &models.ReplyKeyboardRemove{
		RemoveKeyboard: true,
	}
}
