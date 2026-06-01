package service

import (
	"time"

	"prayer-bot/internal/models"
	"prayer-bot/internal/prayer"
)

type PrayerService struct{}

func NewPrayerService() *PrayerService {
	return &PrayerService{}
}

func (s *PrayerService) GetTodaySchedule(
	latitude float64,
	longitude float64,
) (*models.PrayerSchedule, error) {

	wib, _ := time.LoadLocation("Asia/Jakarta")

	now := time.Now().In(wib)

	schedule, err := prayer.CalculatePrayerTimes(
		latitude,
		longitude,
		now,
	)

	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *PrayerService) GetScheduleByDate(
	latitude float64,
	longitude float64,
	date time.Time,
) (*models.PrayerSchedule, error) {

	schedule, err := prayer.CalculatePrayerTimes(
		latitude,
		longitude,
		date,
	)

	if err != nil {
		return nil, err
	}

	return schedule, nil
}
