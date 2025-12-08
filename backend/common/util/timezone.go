package util

import (
	"time"
	_ "time/tzdata"
)

var (
	IST *time.Location
	UTC *time.Location
)

func init() {
	var err error
	IST, err = time.LoadLocation("Asia/Kolkata")
	if err != nil {
		panic("Failed to load Asia/Kolkata timezone: " + err.Error())
	}
	
	UTC = time.UTC
}

func ConvertISTToUTC(istTime time.Time) time.Time {
	if istTime.Location() == time.UTC {
		return time.Date(
			istTime.Year(), istTime.Month(), istTime.Day(),
			istTime.Hour(), istTime.Minute(), istTime.Second(), istTime.Nanosecond(),
			IST,
		).UTC()
	}
	return istTime.In(UTC)
}

func ConvertUTCToIST(utcTime time.Time) time.Time {
	return utcTime.In(IST)
}

func NowUTC() time.Time {
	return time.Now().UTC()
}

func NowIST() time.Time {
	return time.Now().In(IST)
}

func ParseISTTime(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, IST)
}

func ParseISOTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
} 