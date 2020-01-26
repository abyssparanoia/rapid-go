package log

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/pkg/errcode"
	"github.com/abyssparanoia/rapid-go/internal/pkg/util"
)

// Logger ... logger model
type Logger struct {
	Writer            Writer
	MinOutSeverity    Severity
	MaxOuttedSeverity Severity
	TraceID           string
	ResponseStatus    int
	ApplicationLogs   []*EntryChild
}

// IsLogging ... check log level
func (l *Logger) IsLogging(severity Severity) bool {
	return l.MinOutSeverity <= severity
}

// SetOuttedSeverity ... set max log level
func (l *Logger) SetOuttedSeverity(severity Severity) {
	if l.MaxOuttedSeverity < severity {
		l.MaxOuttedSeverity = severity
	}
}

// AddApplicationLog ... add application log
func (l *Logger) AddApplicationLog(severity Severity, file string, line int64, function string, msg string, at time.Time) {
	dst := &EntryChild{
		Severity: severity.String(),
		Message:  fmt.Sprintf("%s:%d [%s] %s", file, line, function, msg),
		Time:     Time(at),
	}
	l.ApplicationLogs = append(l.ApplicationLogs, dst)
}

// WriteRequest ... write request log
func (l *Logger) WriteRequest(r *http.Request, at time.Time, dr time.Duration) {
	l.Writer.Request(
		l.MaxOuttedSeverity,
		l.TraceID,
		l.ApplicationLogs,
		r,
		l.ResponseStatus,
		at,
		dr,
	)
}

// SetResponseStatus ... set response status
func SetResponseStatus(ctx context.Context, status int) {
	logger := GetLogger(ctx)
	if logger != nil {
		logger.ResponseStatus = status
	}
}

// NewLogger ... make logger
func NewLogger(writer Writer, minSeverity Severity, traceID string) *Logger {
	return &Logger{
		Writer:            writer,
		MinOutSeverity:    minSeverity,
		MaxOuttedSeverity: SeverityDebug,
		TraceID:           traceID,
	}
}

// Debugf ... output debug log
func Debugf(ctx context.Context, format string, args ...interface{}) {
	severity := SeverityDebug
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := fmt.Sprintf(format, args...)
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Debugm ... output debug log with define message
func Debugm(ctx context.Context, method string, err error) {
	severity := SeverityDebug
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := fmt.Sprintf("%s: %s", method, err.Error())
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Debuge ... output debug log and make error
func Debuge(ctx context.Context, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityDebug
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return err
}

// Debugc ... output debug log with http status code
func Debugc(ctx context.Context, code int, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityDebug
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return errcode.Set(err, code)
}

// Infof ... output info log
func Infof(ctx context.Context, format string, args ...interface{}) {
	severity := SeverityInfo
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := fmt.Sprintf(format, args...)
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Infom ... output info log with message
func Infom(ctx context.Context, method string, err error) {
	severity := SeverityInfo
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := fmt.Sprintf("%s: %s", method, err.Error())
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Infoe ... output info log and make error
func Infoe(ctx context.Context, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityInfo
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return err
}

// Infoc ... output info log with http status code
func Infoc(ctx context.Context, code int, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityInfo
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return errcode.Set(err, code)
}

// Warningf ... output warning log
func Warningf(ctx context.Context, format string, args ...interface{}) {
	severity := SeverityWarning
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := fmt.Sprintf(format, args...)
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Warningm ... output warning log with message
func Warningm(ctx context.Context, method string, err error) {
	severity := SeverityWarning
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := fmt.Sprintf("%s: %s", method, err.Error())
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Warninge ... output warning log and make error
func Warninge(ctx context.Context, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityWarning
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return err
}

// Warningc ... output warning log with http status code
func Warningc(ctx context.Context, code int, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityWarning
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return errcode.Set(err, code)
}

// Errorf ... output error log
func Errorf(ctx context.Context, format string, args ...interface{}) {
	severity := SeverityError
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := fmt.Sprintf(format, args...)
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Errorm ... output error log with message
func Errorm(ctx context.Context, method string, err error) {
	severity := SeverityError
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := fmt.Sprintf("%s: %s", method, err.Error())
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Errore ... output error log and make error
func Errore(ctx context.Context, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityError
	logger := GetLogger(ctx)
	msg := err.Error()
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return err
}

// Errorc ... output error log with http status code
func Errorc(ctx context.Context, code int, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityError
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return errcode.Set(err, code)
}

// Criticalf ... output critical log
func Criticalf(ctx context.Context, format string, args ...interface{}) {
	severity := SeverityCritical
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := fmt.Sprintf(format, args...)
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Criticalm ... output critical log with message
func Criticalm(ctx context.Context, method string, err error) {
	severity := SeverityCritical
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := fmt.Sprintf("%s: %s", method, err.Error())
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Criticale ... output critical log and make error
func Criticale(ctx context.Context, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityCritical
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return err
}

// Criticalc ... output critical log with http status code
func Criticalc(ctx context.Context, code int, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityCritical
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := util.TimeNow()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return errcode.Set(err, code)
}

// Panic ... handle panic
func Panic(ctx context.Context, rcvr interface{}) string {
	traces := []string{}
	for depth := 0; ; depth++ {
		if depth < 2 {
			continue
		}
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		trace := fmt.Sprintf("%02d: %v:%d", depth-1, file, line)
		traces = append(traces, trace)
	}
	msg := fmt.Sprintf("panic!! %v\n%s", rcvr, strings.Join(traces, "\n"))
	Criticalf(ctx, msg)
	return msg
}

func getFileLine() (string, int64, string) {
	if pt, file, line, ok := runtime.Caller(2); ok {
		parts := strings.Split(file, "/")
		length := len(parts)
		file := fmt.Sprintf("%s/%s", parts[length-2], parts[length-1])

		fParts := strings.Split(runtime.FuncForPC(pt).Name(), ".")
		fLength := len(fParts)
		return file, int64(line), fParts[fLength-1]
	}
	return "", 0, ""
}
