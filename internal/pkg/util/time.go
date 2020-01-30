package util

import (
	"time"
)

// TimeNow ... get recent time of JST
func TimeNow() time.Time {
	return time.Now().In(timeZoneJST())
}

// TimeJST ... get time of JST from Time
func TimeJST(t time.Time) time.Time {
	return t.In(timeZoneJST())
}

// TimeUnix ... get time of JST from unix timestamp
func TimeUnix(u int64) time.Time {
	return time.Unix(u, 0).In(timeZoneJST())
}

func timeZoneJST() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

// TimeNowUnixMill ... get recent time of unixtime mill seccond
func TimeNowUnixMill() int64 {
	return time.Now().In(timeZoneJST()).UnixNano() / int64(time.Millisecond)
}
