package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationService struct {
	db *pgxpool.Pool
}

func NewNotificationService(db *pgxpool.Pool) *NotificationService {
	return &NotificationService{
		db: db,
	}
}

// AlreadySent cek apakah user sudah menerima notifikasi untuk waktu sholat hari ini
func (s *NotificationService) AlreadySent(
	ctx context.Context,
	chatID int64,
	prayerName string,
	date time.Time,
) (bool, error) {

	query := `
		SELECT 1
		FROM prayer_notifications
		WHERE chat_id = $1
		AND prayer_name = $2
		AND prayer_date = $3
		LIMIT 1
	`

	var exists int
	err := s.db.QueryRow(ctx, query, chatID, prayerName, date.Format("2006-01-02")).Scan(&exists)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}
		return false, fmt.Errorf("check already sent: %w", err)
	}

	return true, nil
}

// MarkAsSent catat bahwa user sudah menerima notifikasi untuk waktu sholat hari ini
func (s *NotificationService) MarkAsSent(
	ctx context.Context,
	chatID int64,
	prayerName string,
	date time.Time,
) error {

	query := `
		INSERT INTO prayer_notifications (
			chat_id,
			prayer_name,
			prayer_date
		)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`

	_, err := s.db.Exec(ctx, query, chatID, prayerName, date.Format("2006-01-02"))

	if err != nil {
		return fmt.Errorf("mark as sent: %w", err)
	}

	return nil
}
