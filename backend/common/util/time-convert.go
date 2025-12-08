package util

import (
	"errors"
	"time"
)

const (
	inputTimeLayout = "2006-01-02T15:04:05" // Expecting IST in this format
)

var (
	istLocation, _ = time.LoadLocation("Asia/Kolkata")
)

// ConvertISTStringToUTC parses "YYYY-MM-DDTHH:MM:SS" in IST and returns UTC time
func ConvertISTStringToUTC(ts string) (time.Time, error) {
	t, err := time.ParseInLocation(inputTimeLayout, ts, istLocation)
	if err != nil {
		return time.Time{}, errors.New("invalid scheduled_at timestamp format: must be YYYY-MM-DDTHH:MM:SS in IST")
	}
	return t.UTC(), nil
}
