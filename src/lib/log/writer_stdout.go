package log

import (
	"fmt"
	"net/http"
	"time"
)

type writerStdout struct {
	TimeFormat string
}

func (w *writerStdout) Request(
	severity Severity,
	traceID string,
	applicationLogs []*EntryChild,
	r *http.Request,
	status int,
	at time.Time,
	dr time.Duration) {
	u := *r.URL
	u.Fragment = ""
	date := at.Format(w.TimeFormat)
	fmt.Printf("%s \"%s %s\" %d %dms\n", date, r.Method, u.RequestURI(), status, dr/1000000)
}

func (w *writerStdout) Application(
	severity Severity,
	traceID string,
	msg string,
	file string,
	line int64,
	function string,
	at time.Time) {
	date := at.Format(w.TimeFormat)
	fmt.Printf("%s [%s] %s:%d [%s] %s\n", date, severity.String(), file, line, function, msg)
}

// NewWriterStdout ... get writer for stdout
func NewWriterStdout() Writer {
	return &writerStdout{
		TimeFormat: "2006-01-02 15:04:05.000",
	}
}
