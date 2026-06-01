package prayer

import (
	"time"

	adhan "github.com/hablullah/go-prayer"
)

type PrayerSchedule struct {
	Fajr    time.Time
	Dhuhr   time.Time
	Asr     time.Time
	Maghrib time.Time
	Isha    time.Time
}

func CalculatePrayerTimes(
	latitude float64,
	longitude float64,
	date time.Time,
) (*PrayerSchedule, error) {

	coordinates := adhan.Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}

	params := adhan.GetMuslimWorldLeagueConfiguration()

	// Indonesia mayoritas Syafi'i
	params.Madhab = adhan.Shafi

	// MABIMS/Kemenag approximation
	params.FajrAngle = 20.0
	params.IshaAngle = 18.0

	prayerTimes := adhan.NewPrayerTimes(
		coordinates,
		date,
		params,
	)

	schedule := &PrayerSchedule{
		Fajr:    prayerTimes.Fajr,
		Dhuhr:   prayerTimes.Dhuhr,
		Asr:     prayerTimes.Asr,
		Maghrib: prayerTimes.Maghrib,
		Isha:    prayerTimes.Isha,
	}

	return schedule, nil
}
