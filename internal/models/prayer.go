package models

import "time"

type PrayerSchedule struct {
	Fajr    time.Time
	Dhuhr   time.Time
	Asr     time.Time
	Maghrib time.Time
	Isha    time.Time
}

func (p PrayerSchedule) PrayerMap() map[string]time.Time {
	return map[string]time.Time{
		"fajr":    p.Fajr,
		"dhuhr":   p.Dhuhr,
		"asr":     p.Asr,
		"maghrib": p.Maghrib,
		"isha":    p.Isha,
	}
}
