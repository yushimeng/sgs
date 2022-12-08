package log

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	logger *Logger
)

type Logger struct {
	lvl         LogLevel
	rec         LogRecordChan // 初始化时，设置大小
	file        *os.File
	filename    string
	out_chan    chan string
	close_chan  chan int
	rotate_chan chan bool
}

// func (logger *Logger) FormatToString() (str string) {
// 	// buf := bytes.NewBuffer([]byte, 0, 64)
// 	now := time.Now()
// 	month, day, year := now.Month(), now.Day(), now.Year()
// 	hour, minute, second := now.Hour(), now.Minute(), now.Second()
// 	zone, _ := now.Zone()
// 	out := bytes.NewBuffer(make([]byte, 0, 64))
// 	out.WriteString("---------------------------\n")
// 	out.WriteString("Logger Info report: \n")
// 	fmt.Fprintf(out, "Time: %04d-%02d-%02d %02d:%02d:%02d %s\n", year, month, day, hour, minute, second, zone)
// 	fmt.Fprintf(out, "Log level: %s\n", logger.lvl.FormatToString())
// 	out.WriteString("---------------------------\n")

// 	return out.String()
// }

func init() {
	logger = &Logger{
		lvl:         INFO,
		rec:         make(chan *LogRecord, 100),
		out_chan:    make(chan string),
		close_chan:  make(chan int),
		rotate_chan: make(chan bool),
	}
	go logger.consumer()
}

func SetLogLevel(lvl LogLevel) {
	logger.lvl = lvl
}

func SetRotateDaily() {
	go logger.writeRotate()
}

func SetOutput(filename string) error {
	var err error

	logger.out_chan <- filename

	res, ok := <-logger.out_chan
	if !ok {
		return errors.New("recv set output result error")
	}
	if res != filename {
		return errors.New("OpenFile error")
	}

	return err
}

func Debug(arg0 interface{}, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	if lvl < logger.lvl {
		return
	}

	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		logger.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		logger.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		logger.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func Trace(arg0 interface{}, args ...interface{}) {
	const (
		lvl = TRACE
	)
	if lvl < logger.lvl {
		return
	}

	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		logger.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		logger.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		logger.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func Info(arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)
	if lvl < logger.lvl {
		return
	}

	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		logger.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		logger.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		logger.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func Warn(arg0 interface{}, args ...interface{}) {
	const (
		lvl = WARN
	)
	if lvl < logger.lvl {
		return
	}

	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		logger.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		logger.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		logger.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func Error(arg0 interface{}, args ...interface{}) {
	const (
		lvl = ERROR
	)
	if lvl < logger.lvl {
		return
	}

	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		logger.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		logger.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		logger.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func Fatal(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FATAL
	)
	if lvl < logger.lvl {
		return
	}

	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		logger.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		logger.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		logger.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func (logger *Logger) intLogf(lvl LogLevel, format string, args ...interface{}) {
	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
		npos := strings.LastIndex(src, "/")
		if npos > 0 {
			src = string([]byte(src[npos+1:]))
		}
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Source:  src,
		Message: msg,
	}

	// Dispatch the logs
	logger.Write(rec)
}

// Send a closure log message internally
func (log Logger) intLogc(lvl LogLevel, closure func() string) {
	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
		npos := strings.LastIndex(src, "/")
		if npos > 0 {
			src = string([]byte(src[npos+1:]))
		}
	}

	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Source:  src,
		Message: closure(),
	}

	// Dispatch the logs
	logger.Write(rec)
}

func (logger *Logger) Write(record *LogRecord) {
	logger.rec <- record
}

func (logger *Logger) writeRotate() {
	date := time.Now().Format("2006-01-02")
	for {
		time.Sleep(1 * time.Minute)
		curDate := time.Now().Format("2006-01-02")
		if curDate != date {
			date = curDate
			logger.rotate_chan <- true
		}
	}
}

func Close() {
	logger.close_chan <- 1
}
