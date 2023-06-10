package main

import (
	"fmt"
	"time"
)

func convertTimezones(date time.Time) (time.Time, time.Time) {
	newYorkLocation, _ := time.LoadLocation("America/New_York")
	bogotaLocation, _ := time.LoadLocation("America/Bogota")

	newYorkTime := date.In(newYorkLocation)
	bogotaTime := date.In(bogotaLocation)

	return newYorkTime, bogotaTime
}

func main() {
	// Example usage with a date in UTC
	utcDate := time.Date(2023, time.March, 12, 7, 0, 0, 0, time.UTC)

	newYorkTime, bogotaTime := convertTimezones(utcDate)

	fmt.Println("UTC Date:", utcDate)
	fmt.Println("New York Time:", newYorkTime)
	fmt.Println("Bogota Time:", bogotaTime)

	// Calculate the time difference between New York and Bogota
	timeDiff := bogotaTime.Hour() - newYorkTime.Hour()
	fmt.Println("Time difference between New York and Bogota:", timeDiff)
}
