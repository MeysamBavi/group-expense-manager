package model

import (
	"fmt"
	ptime "github.com/yaa110/go-persian-calendar"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var layouts = [...]string{
	"2006/1/2 15:4",
	"15:4 2006/1/2",
	"2006/1/2",
	"2006-1-2",
	"2006 1 2",

	// possible auto format results:
	"1/2/2006 15:4",
	"1/2/06 15:4",
	"2/1/2006 15:4",
	"2/1/06 15:4",
	"2-1-2006",
	"2-1-06",
	"1-2-2006",
	"1-2-06",
	"1/2/2006",
	"1/2/06",
	"2/1/2006",
	"2/1/06",
}

type Time interface {
	fmt.Stringer
}

func ParseTime(value string) (Time, error) {
	g, errG := parseGregorian(value)
	if errG == nil {
		year := g.Year()
		if year >= 2000 {
			return g, nil
		}
	}

	p, errP := parsePersian(value)
	if errP == nil {
		return p, nil
	}

	if errG == nil {
		return g, nil
	}

	return nil, fmt.Errorf("time value %q cannot be parsed:\n\t%v\n\t%v", value, errG, errP)
}

type gregorian struct {
	time.Time
}

func (g *gregorian) String() string {
	if g == nil || g.Time.IsZero() {
		return ""
	}

	return g.Format("2006/01/02 15:04")
}

func parseGregorian(value string) (*gregorian, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return &gregorian{}, nil
	}
	var firstError error
	for _, layout := range layouts {
		if result, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &gregorian{result}, nil
		} else if firstError == nil {
			firstError = err
		}
	}

	return nil, fmt.Errorf("cannot not parse %q as gregorian date: %w", value, firstError)
}

func TimeOfGregorian(t time.Time) Time {
	return &gregorian{t}
}

type persian struct {
	ptime.Time
}

func (p *persian) String() string {
	return p.Format("yyyy/MM/dd HH:mm")
}

const (
	dayRE    = "(?P<day>\\d{1,2})"
	monthRE  = "(?P<month>\\d{1,2})"
	yearRE   = "(?P<year>\\d{4})"
	hourRE   = "(?P<hour>\\d{1,2})"
	minuteRE = "(?P<minute>\\d{1,2})"
)

var persianLayouts = []*regexp.Regexp{
	regexp.MustCompile(fmt.Sprintf("^%s/%s/%s %s:%s$", yearRE, monthRE, dayRE, hourRE, minuteRE)),
	regexp.MustCompile(fmt.Sprintf("^%s:%s %s/%s/%s$", hourRE, minuteRE, yearRE, monthRE, dayRE)),
	regexp.MustCompile(fmt.Sprintf("^%s/%s/%s %s:%s$", dayRE, monthRE, yearRE, hourRE, minuteRE)),
	regexp.MustCompile(fmt.Sprintf("^%s:%s %s/%s/%s$", hourRE, minuteRE, dayRE, monthRE, yearRE)),
	regexp.MustCompile(fmt.Sprintf("^%s/%s/%s$", yearRE, monthRE, dayRE)),
	regexp.MustCompile(fmt.Sprintf("^%s-%s-%s$", yearRE, monthRE, dayRE)),
	regexp.MustCompile(fmt.Sprintf("^%s %s %s$", yearRE, monthRE, dayRE)),
	regexp.MustCompile(fmt.Sprintf("^%s/%s/%s$", dayRE, monthRE, yearRE)),
	regexp.MustCompile(fmt.Sprintf("^%s-%s-%s$", dayRE, monthRE, yearRE)),
	regexp.MustCompile(fmt.Sprintf("^%s %s %s$", dayRE, monthRE, yearRE)),
}

func parsePersian(str string) (*persian, error) {
	str = strings.TrimSpace(str)
	var mainError *persianParseError
	for _, re := range persianLayouts {
		result, err := parsePersianWithLayout(str, re)
		if err == nil {
			return result, nil
		}

		if mainError == nil || err.priority > mainError.priority {
			mainError = err
		}
	}

	return nil, fmt.Errorf("cannot not parse %q as persian date: %w", str, mainError.err)
}

type persianParseError struct {
	err      error
	priority int
}

func parsePersianWithLayout(str string, re *regexp.Regexp) (*persian, *persianParseError) {
	var day, month, year, hour, minute int
	type parsableField struct {
		value      *int
		validRange [2]int
	}
	m := map[string]*parsableField{
		"day":    {value: &day, validRange: [2]int{1, 31}},
		"month":  {value: &month, validRange: [2]int{1, 12}},
		"year":   {value: &year, validRange: [2]int{100, 9999}},
		"hour":   {value: &hour, validRange: [2]int{0, 23}},
		"minute": {value: &minute, validRange: [2]int{0, 59}},
	}
	names := re.SubexpNames()
	result := re.FindStringSubmatch(str)
	if result == nil {
		return nil, &persianParseError{fmt.Errorf("%q does not match %q", str, re.String()), 0}
	}
	for i, name := range names {
		if name == "" {
			continue
		}
		value, err := strconv.Atoi(result[i])
		if err != nil {
			return nil, &persianParseError{fmt.Errorf("cannot not parse %q as int: %w", result[i], err), 1}
		}
		theRange := m[name].validRange
		if value > theRange[1] || value < theRange[0] {
			return nil, &persianParseError{fmt.Errorf("value %d is not in range %v", value, theRange), 2}
		}
		*m[name].value = value
	}

	return &persian{ptime.Date(year, ptime.Month(month), day, hour, minute, 0, 0, time.Local)}, nil
}
