package utils

import "time"

func GetEuropeTime() time.Time {
	return time.Now().UTC().Add(time.Hour * 3)
}
