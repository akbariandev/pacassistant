package prettytime

import (
	"fmt"
)

type TimeScale string

const (
	SECONDS TimeScale = "secs"
	MINUTES TimeScale = "min"
	HOURS             = "hours"
	DAYS              = "days"
)

func ElapsedTime(time uint) string {
	var s TimeScale
	var t uint

	switch {
	case time < 60:
		t = time
		s = SECONDS
	case time >= 60 && time < 3600:
		t = time / 60
		s = MINUTES
	case time >= 3600 && time < 86400:
		t = time / 3600
		s = HOURS
	case time >= 86400:
		t = time / 86400
		s = DAYS
	}

	return fmt.Sprintf("%d %s ago", t, s)
}

func GetDurationComponents(duration int) (days, hours, mins, secs int) {
	days = duration / (24 * 3600)
	duration %= 24 * 3600

	hours = duration / 3600
	duration %= 3600

	mins = duration / 60
	duration %= 60

	secs = duration

	return days, hours, mins, secs
}
