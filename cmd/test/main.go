package main

import (
	"fmt"
	"time"

	"prayer-bot/internal/prayer"
)

func main() {

	wib, _ := time.LoadLocation(
		"Asia/Jakarta",
	)

	schedule, err := prayer.CalculatePrayerTimes(
		-0.357229,
		102.314851,
		time.Now().In(wib),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Subuh   :", schedule.Fajr.Format("15:04"))
	fmt.Println("Dzuhur  :", schedule.Dhuhr.Format("15:04"))
	fmt.Println("Ashar   :", schedule.Asr.Format("15:04"))
	fmt.Println("Maghrib :", schedule.Maghrib.Format("15:04"))
	fmt.Println("Isya    :", schedule.Isha.Format("15:04"))
}
