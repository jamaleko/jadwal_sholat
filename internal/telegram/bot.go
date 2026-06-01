package telegram

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Bot struct {
	client *bot.Bot
	db     *pgxpool.Pool
}

func NewBot(
	token string,
	db *pgxpool.Pool,
) (*Bot, error) {

	client, err := bot.New(token)
	if err != nil {
		return nil, err
	}

	telegramBot := &Bot{
		client: client,
		db:     db,
	}
	
	telegramBot.registerHandlers()
	
	return telegramBot, nil
}

func (b *Bot) Start() {
	log.Println("Telegram bot started")

	b.client.Start(context.Background())
}

func (b *Bot) Client() *bot.Bot {
	return b.client
}

func (b *Bot) DB() *pgxpool.Pool {
	return b.db
}
