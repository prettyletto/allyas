package datetime

import "time"

const HumanLayout = "01/02/2006 15:04"

func Format(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Local().Format(HumanLayout)
}
