package log

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	TRACE
	INFO
	WARN
	ERROR
	FATAL
)

var (
	LogLevelStr = [...]string{"DEBUG", "TRACE", "INFO", "WARN", "ERROR", "FATAL"}
)

// A LogRecord contains all of the pertinent information for each message
type LogRecord struct {
	Level   LogLevel  // The log level
	Created time.Time // The time at which the log message was created (nanoseconds)
	Source  string    // The message source
	Message string    // The log message
}

type LogRecordChan chan *LogRecord

func (rec LogRecord) Format() string {
	month, day, year := rec.Created.Month(), rec.Created.Day(), rec.Created.Year()
	hour, minute, second := rec.Created.Hour(), rec.Created.Minute(), rec.Created.Second()
	nano := rec.Created.UnixNano()
	zone, _ := rec.Created.Zone()

	msg := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d,%03d %s [%s][%s] %s\n",
		year, month, day, hour, minute, second, (nano%1e9)/1e6, zone,
		LogLevelStr[int(rec.Level)], rec.Source, rec.Message)
	return msg
}
