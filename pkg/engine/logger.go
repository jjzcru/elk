package engine

import (
	"fmt"
	"github.com/jjzcru/elk/pkg/file"
	"io"
	"os"
	"strings"
	"time"
)

// Logger is used by the engine to store the output
type Logger struct {
	StderrWriter io.Writer
	StdoutWriter io.Writer
	StdinReader  io.Reader
}

// DefaultLogger is the standard output for a logger
func DefaultLogger() Logger {
	return Logger{
		StdoutWriter: os.Stdout,
		StderrWriter: os.Stderr,
		StdinReader:  os.Stdin,
	}
}

// TimeStampLogger receives any logger and appends a timestamp in a specific format
func TimeStampLogger(logger Logger, format string) (Logger, error) {
	formatter, err := getTimeStampFormatter(format)
	if err != nil {
		return Logger{}, err
	}

	timeStampLogger := Logger{
		StderrWriter: TimeStampWriter{
			writer:    logger.StderrWriter,
			TimeStamp: formatter,
		},
		StdoutWriter: TimeStampWriter{
			writer:    logger.StdoutWriter,
			TimeStamp: formatter,
		},
		StdinReader: logger.StdinReader,
	}

	return timeStampLogger, nil
}

// TimeStampWriter attach a timestamp to each log
type TimeStampWriter struct {
	writer    io.Writer
	TimeStamp func() string
}

func (t TimeStampWriter) GetDateFormat(format string) (string, error) {
	switch format {
	case "ANSIC":
		fallthrough
	case "ansic":
		return time.ANSIC, nil
	case "UnixDate":
		fallthrough
	case "unixdate":
		return time.UnixDate, nil
	case "rubydate":
		fallthrough
	case "RubyDate":
		return time.RubyDate, nil
	case "RFC822":
		return time.RFC822, nil
	case "RFC822Z":
		return time.RFC822Z, nil
	case "RFC850":
		return time.RFC850, nil
	case "RFC1123":
		return time.RFC1123, nil
	case "RFC1123Z":
		return time.RFC1123Z, nil
	case "RFC3339":
		return time.RFC3339, nil
	case "RFC3339Nano":
		return time.RFC3339Nano, nil
	case "kitchen":
		fallthrough
	case "Kitchen":
		return time.Kitchen, nil
	default:
		return "", fmt.Errorf("%s is an invalid timestamp format", format)
	}
}

// Writes timestamp to a writer
func (t TimeStampWriter) Write(p []byte) (int, error) {
	var err error
	if t.writer != nil {
		content := string(p)
		breakLine := file.BreakLine

		timestamp := t.TimeStamp()
		timeStampPrefix := fmt.Sprintf("%s%s | ", file.BreakLine, timestamp)
		var inputs []string
		for _, input := range strings.SplitN(content, breakLine, -1) {
			if len(input) > 0 {
				inputs = append(inputs, input)
			}
		}

		response := timeStampPrefix + strings.Join(inputs, timeStampPrefix)

		_, err = t.writer.Write([]byte(response))
	}
	return len(p), err
}

func getTimeStampFormatter(format string) (func() string, error) {
	var formatter func() string
	switch format {
	case time.ANSIC:
		fallthrough
	case time.UnixDate:
		fallthrough
	case time.RubyDate:
		fallthrough
	case time.RFC822:
		fallthrough
	case time.RFC822Z:
		fallthrough
	case time.RFC850:
		fallthrough
	case time.RFC1123:
		fallthrough
	case time.RFC1123Z:
		fallthrough
	case time.RFC3339:
		fallthrough
	case time.RFC3339Nano:
		fallthrough
	case time.Kitchen:
		formatter = func() string {
			return time.Now().Format(format)
		}
	default:
		return nil, fmt.Errorf("invalid date format")
	}

	return formatter, nil
}
