package prayer

import (
	"time"

	goprayer "github.com/hablullah/go-prayer"

	"prayer-bot/internal/models"
)

func CalculatePrayerTimes(
	latitude float64,
	longitude float64,
	date time.Time,
) (*models.PrayerSchedule, error) {

	location := date.Location()

	schedules, err := goprayer.Calculate(
		goprayer.Config{
			Latitude:           latitude,
			Longitude:          longitude,
			Timezone:           location,
			TwilightConvention: goprayer.Kemenag(),
			AsrConvention:      goprayer.Shafii,
			PreciseToSeconds:   false,
		},
		date.Year(),
	)

	if err != nil {
		return nil, err
	}

	day := date.YearDay()

	schedule := schedules[day-1]
	loc := date.Location()

	result := &models.PrayerSchedule{
		Fajr:    schedule.Fajr.In(loc),
		Dhuhr:   schedule.Zuhr.In(loc),
		Asr:     schedule.Asr.In(loc),
		Maghrib: schedule.Maghrib.In(loc),
		Isha:    schedule.Isha.In(loc),
	}

	return result, nil
}
