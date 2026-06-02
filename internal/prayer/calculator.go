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
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
	 return nil, err
	}

	result := &models.PrayerSchedule{
	 Fajr:    schedule.Fajr.In(location),
	 Dhuhr:   schedule.Zuhr.In(location),
	 Asr:     schedule.Asr.In(location),
	 Maghrib: schedule.Maghrib.In(location),
	 Isha:    schedule.Isha.In(location),
	}

	return result, nil
}
