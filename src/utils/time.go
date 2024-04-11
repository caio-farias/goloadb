package utils

import (
	"math"
	"time"
)

const (
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	BR_Location = "America/Sao_Paulo"
)

func getTime(location_name string, format string) (string, error) {
	location, error := time.LoadLocation(location_name)
	if error != nil {
		return "", error
	}

	return time.Now().In(location).Format(format), nil
}

func GetTimeHere() string {
	time, _ := getTime(BR_Location, RFC1123)
	return time
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64) float64 {
	output := math.Pow(10, float64(4))
	return float64(round(num*output)) / output
}

func GetTimeInSeconds(sec int) time.Duration {
	return time.Duration(sec) * time.Second
}

func GetTimeInMilliseconds(mills int) time.Duration {
	return time.Duration(mills) * time.Millisecond

}

func GetTimeDiffInMicrosec(start time.Time) float64 {
	return ToFixed(float64(time.Since(start).Microseconds()))
}

func GetTimeDiffInMillisec(start time.Time) float64 {
	return ToFixed(float64(time.Since(start).Milliseconds()))
}
