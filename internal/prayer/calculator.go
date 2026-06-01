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
		Fajr:    schedule.Fajr,
		Dhuhr:   schedule.Zuhr,
		Asr:     schedule.Asr,
		Maghrib: schedule.Maghrib,
		Isha:    schedule.Isha,
	}

	return result, nil
}
