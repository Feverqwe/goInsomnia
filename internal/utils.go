package internal

import "strconv"

func formatMinutes(minutes int) string {
	if minutes < 60 {
		return strconv.Itoa(minutes) + " minutes"
	}
	return strconv.Itoa(minutes/60) + " hours"
}
