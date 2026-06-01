package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"prayer-bot/internal/models"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(
	db *pgxpool.Pool,
) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) SaveLocation(
	ctx context.Context,
	chatID int64,
	latitude float64,
	longitude float64,
) error {

	query := `
		INSERT INTO users (
			chat_id,
			latitude,
			longitude,
			active
		)
		VALUES (
			$1,
			$2,
			$3,
			true
		)
		ON CONFLICT (chat_id)
		DO UPDATE SET
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			active = true;
	`

	_, err := s.db.Exec(
		ctx,
		query,
		chatID,
		latitude,
		longitude,
	)

	if err != nil {
		return fmt.Errorf(
			"save user location: %w",
			err,
		)
	}

	return nil
}

func (s *UserService) GetByChatID(
	ctx context.Context,
	chatID int64,
) (*models.User, error) {

	query := `
		SELECT
			chat_id,
			latitude,
			longitude,
			active,
			created_at
		FROM users
		WHERE chat_id = $1
	`

	var user models.User

	err := s.db.QueryRow(
		ctx,
		query,
		chatID,
	).Scan(
		&user.ChatID,
		&user.Latitude,
		&user.Longitude,
		&user.Active,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (s *UserService) Enable(
	ctx context.Context,
	chatID int64,
) error {

	query := `
		UPDATE users
		SET active = true
		WHERE chat_id = $1
	`

	_, err := s.db.Exec(
		ctx,
		query,
		chatID,
	)

	if err != nil {
		return fmt.Errorf(
			"enable user: %w",
			err,
		)
	}

	return nil
}
func (s *UserService) GetActiveUsers(
	ctx context.Context,
) ([]models.User, error) {

	query := `
		SELECT
			chat_id,
			latitude,
			longitude,
			active,
			created_at
		FROM users
		WHERE active = true
	`

	rows, err := s.db.Query(
		ctx,
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {

		var user models.User

		err := rows.Scan(
			&user.ChatID,
			&user.Latitude,
			&user.Longitude,
			&user.Active,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(
			users,
			user,
		)
	}

	return users, nil
}

func (s *UserService) Disable(
	ctx context.Context,
	chatID int64,
) error {

	query := `
		UPDATE users
		SET active = false
		WHERE chat_id = $1
	`

	_, err := s.db.Exec(
		ctx,
		query,
		chatID,
	)

	if err != nil {
		return fmt.Errorf(
			"disable user: %w",
			err,
		)
	}

	return nil
}
