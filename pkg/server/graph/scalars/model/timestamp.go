package model

import (
	"errors"
	"io"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalTimestamp marshall time stamp content as RFC3339
func MarshalTimestamp(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, t.Format(time.RFC3339))
	})
}

// UnmarshalTimestamp transform value to Timestamps
func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(int64); ok {
		return time.Unix(tmpStr, 0), nil
	}

	if tmpStr, ok := v.(string); ok {
		validTimeFormats := []string{
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			time.RFC822,
			time.RFC822Z,
			time.RFC850,
			time.RFC1123,
			time.RFC1123Z,
			time.RFC3339,
			time.RFC3339Nano,
			time.Kitchen,
		}

		for _, layout := range validTimeFormats {
			switch layout {
			case time.Kitchen:
				fallthrough
			case time.RFC3339:
				fallthrough
			case time.RFC3339Nano:
				tmpStr = strings.Replace(tmpStr, " ", "", -1)
			}
			deadlineTime, err := time.Parse(layout, tmpStr)
			if err == nil {
				if layout == time.Kitchen {
					now := time.Now()
					deadlineTime = time.Date(now.Year(),
						now.Month(),
						now.Day(),
						deadlineTime.Hour(),
						deadlineTime.Minute(),
						0,
						0,
						now.Location())

					// If time is before now i refer to that time but the next day
					if deadlineTime.Before(now) {
						deadlineTime = deadlineTime.Add(24 * time.Hour)
					}
				}
				return deadlineTime, nil
			}
		}

		return time.Time{}, errors.New("invalid date format")
	}

	return time.Time{}, errors.New("time should be a string or a unix timestamp")
}
