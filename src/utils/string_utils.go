package utils

import (
	"log"
	"strings"
	"time"
)

func SanitizeString(str string) string {
	return strings.ReplaceAll(str, "\r", "")
}

func Unpack(s []string, vars ...*string) {
	len := len(vars)
	for i, str := range s {
		if i >= len {
			break
		}
		*vars[i] = str
	}
}

func PrintAndSleep(duration int, message string, vars ...any) {
	log.Printf(message, vars...)
	time.Sleep(GetTimeInSeconds(duration))
}
