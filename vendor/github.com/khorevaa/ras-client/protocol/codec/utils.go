package codec

import "time"

const AgeDelta = 621355968000000

func dateFromTicks(ticks int64) time.Time {
	if ticks > 0 {

		timeT := (ticks - AgeDelta) / 10

		t := time.Unix(0, timeT*int64(time.Millisecond))

		return t

	}
	return time.Time{}
}

func dateToTicks(date time.Time) (ticks int64) {

	if !date.IsZero() {

		ticks = date.UnixNano() / int64(time.Millisecond)

		ticks = ticks*10 + AgeDelta

		return ticks

	}
	return 0
}
