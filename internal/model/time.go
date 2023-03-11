package model

import (
	"fmt"
	"strings"
	"time"
)

var layouts = [...]string{
	"2006/01/02 15:04",
	"2006/01/02 15:4",
	"2006/01/02 15:04:05",
	"2006/01/02 15:4:5",
	"15:04 2006/01/02",
	"15:4 2006/01/02",
	"15:04:05 2006/01/02",
	"15:4:5 2006/01/02",
	"2006/1/2 15:04",
	"2006/1/2 15:4",
	"2006/1/2 15:04:05",
	"2006/1/2 15:4:5",
	"15:04 2006/1/2",
	"15:4 2006/1/2",
	"15:04:05 2006/1/2",
	"15:4:5 2006/1/2",
	"2006/01/02",
	"2006/1/2",
	"2006-01-02",
	"2006-1-2",
	"2006 01 02",
	"2006 1 2",
}

type Time interface {
	fmt.Stringer
}

type gregorian struct {
	time.Time
}

func (g *gregorian) String() string {
	if g == nil || g.Time.IsZero() {
		return ""
	}

	return g.Format(layouts[0])
}

func ParseTime(value string) (Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return &gregorian{}, nil
	}
	for _, layout := range layouts {
		if result, err := time.Parse(layout, value); err == nil {
			return &gregorian{result.Local()}, nil
		}
	}

	return nil, fmt.Errorf("time value %q does not match any layout", value)
}

func TimeOfGregorian(t time.Time) Time {
	return &gregorian{t}
}
