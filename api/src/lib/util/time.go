package util

import (
	"time"
)

// TimeNow ... 現在時刻をJSTのTimeで取得する
func TimeNow() time.Time {
	return time.Now().In(timeZoneJST())
}

// TimeJST ... TimeからJSTのTimeを取得する
func TimeJST(t time.Time) time.Time {
	return t.In(timeZoneJST())
}

// TimeUnix ... UnixTimestampからJSTのTimeを取得する
func TimeUnix(u int64) time.Time {
	return time.Unix(u, 0).In(timeZoneJST())
}

func timeZoneJST() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}
