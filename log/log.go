package log

import (
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
	rec      LogRecordChan // 初始化时，设置大小
	file     *os.File
	filename string
}

func init() {
	logger = &Logger{
		rec: make(chan *LogRecord, 100),
	}
	go logger.output()
}

func SetOutput(filename string) (err error) {
	return logger.SetOutput(filename)
}

func (logger Logger) SetOutput(file_name string) (err error) {
	if file_name == logger.filename {
		return err
	}

	if logger.file != nil {
		logger.file.Close()
	}

	logger.filename = file_name
	if logger.filename != "" {
		fd, err := os.OpenFile(file_name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		logger.file = fd
		fmt.Printf("SetOutput-》设置日志成功 logger%p.file=%p\n", &logger, logger.file)
		fmt.Fprintf(logger.file, "open file[%s]----------------\n", logger.filename)
	}
	return err
}

func Info(arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)
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

// Glog = Logger
func (logger *Logger) output() {
	defer func() {
		if logger.file != nil {
			logger.file.Close()
		}
		close(logger.rec)
	}()
	for {
		select {
		case rec, ok := <-logger.rec:
			if !ok {
				fmt.Println("read rec failed")
				return
			}
			fmt.Printf("init->output:   logger:%p, file:%p\n", logger, logger.file)
			// fmt.Fprint(logger.file, rec.Print())
			if logger.file != nil {
				// Perform the write
				// _, err := fmt.Fprint(logger.file, rec.Format())
				_, err := fmt.Fprintf(logger.file, "%s", rec.Format())
				if err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", logger.filename, err)
					return
				}
			} else {
				fmt.Println(rec.Format())
			}
		}
	}
}
func Close() {

}
