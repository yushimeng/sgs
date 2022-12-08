package log

import (
	"fmt"
	"os"
	"time"
)

func (logger *Logger) consumer() {
	defer func() {
		if logger.file != nil {
			logger.file.Close()
		}
		close(logger.rec)
	}()
	for {
		select {
		case out, ok := <-logger.out_chan:
			if !ok {
				fmt.Println("read out chan failed")
				return
			}
			logger.consumerFile(out)
		case _, ok := <-logger.rotate_chan:
			if !ok {
				fmt.Println("read rotate chan failed")
				return
			}
			logger.consumerRotate()
		default:
			select {
			case rec, ok := <-logger.rec:
				if !ok {
					fmt.Println("read rec failed")
					return
				}
				logger.consumerRecord(rec)

			default:
				select {
				case _, ok := <-logger.close_chan:
					if !ok {
						fmt.Println("failed read close chan")
					}
					return
				default:
				}
			}
		}
	}
}

func (logger *Logger) consumerRecord(rec *LogRecord) (err error) {
	if logger.file != nil {
		// Perform the write
		_, err = fmt.Fprintf(logger.file, "%s", rec.Format())
		if err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s", logger.filename, err)
			return
		}
	} else {
		fmt.Print(rec.Format())
	}
	return
}

func (logger *Logger) consumerFile(out string) (err error) {
	defer func() {
		logger.out_chan <- out
	}()
	if out == logger.filename {
		return err
	}
	if logger.file != nil {
		logger.file.Close()
	}
	logger.filename = out
	if logger.filename != "" {
		fd, err := os.OpenFile(logger.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("failed to open new file", out)
			logger.filename = ""
			return err
		}
		logger.file = fd
		// fmt.Fprintf(logger.file, "%s", logger.FormatToString())
	}
	return err
}

func (logger *Logger) consumerRotate() (err error) {
	if logger.file == nil {
		return err
	}
	logger.file.Close()
	fname := logger.filename + "." + time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	os.Rename(logger.filename, fname)

	fd, err := os.OpenFile(logger.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("failed to open file")
		return err
	}

	logger.file = fd

	return err
}
