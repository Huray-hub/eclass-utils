package assignment

import (
	"errors"

	"strings"
	"time"
)

var dayNamesGR = map[string]string{
	"Δευτέρα":   "Monday",
	"Τρίτη":     "Tuesday",
	"Τετάρτη":   "Wednesday",
	"Πέμπτη":    "Thursday",
	"Παρασκευή": "Friday",
	"Σάββατο":   "Saturday",
	"Κυριακή":   "Sunday",
}

var monthNamesGenitiveGR = map[string]string{
	"Ιανουαρίου":  "January",
	"Φεβρουαρίου": "February",
	"Μαρτίου":     "March",
	"Απριλίου":    "April",
	"Μαΐου":       "May",
	"Ιουνίου":     "June",
	"Ιουλίου":     "July",
	"Αυγούστου":   "August",
	"Σεπτεμβρίου": "September",
	"Οκτωβρίου":   "October",
	"Νοεμβρίου":   "November",
	"Δεκεμβρίου":  "December",
}

var periodsGR = map[string]string{
	"π.μ.": "am",
	"μ.μ.": "pm",
}

func parseTime(dateRaw string, location *time.Location) (*time.Time, error) {
	if strings.Contains(dateRaw, "αύριο") || strings.Contains(dateRaw, "μεθαύριο") ||
		strings.Contains(dateRaw, "σήμερα") {
		return parseNearTime(dateRaw, location)
	}

	return parseNormalDate(dateRaw, location)
}

// parseNearTime parses the following formats:
// "αύριο - 11:59 μ.μ.(απομένουν 1 ημέρα 3 ώρες 8 λεπτά)"
// "μεθαύριο - 11:59 μ.μ.(απομένουν 2 ημέρες 3 λώρες 8 λεπτά)"
func parseNearTime(nearDate string, location *time.Location) (*time.Time, error) {
	dateOnly, err := parseNearDateOnly(nearDate)
	if err != nil {
		return nil, err
	}

	timeOnly, err := parseTimeOnly(nearDate)
	if err != nil {
		return nil, err
	}

	fullTime := time.Date(
		dateOnly.Year(),
		dateOnly.Month(),
		dateOnly.Day(),
		timeOnly.Hour(),
		timeOnly.Minute(),
		timeOnly.Second(),
		0,
		location,
	)

	return &fullTime, nil
}

func parseNearDateOnly(date string) (*time.Time, error) {
	dateOnly := time.Now()

	switch {
	case strings.Contains(date, "αύριο"):
		dateOnly = dateOnly.AddDate(0, 0, 1)
	case strings.Contains(date, "μεθαύριο"):
		dateOnly = dateOnly.AddDate(0, 0, 2)
	default:
		return nil, errors.New("invalid value")
	}

	return &dateOnly, nil
}

func parseTimeOnly(timeRaw string) (*time.Time, error) {
	timeParts := strings.Split(timeRaw, " ")
	periodGR := strings.Split(timeParts[3], "(")[0]

	timeRaw = timeParts[2] + periodsGR[periodGR]

	time, err := time.Parse("15:04pm", timeRaw)
	if err != nil {
		return nil, err
	}

	return &time, nil

}

// parseNormalDate will parse
func parseNormalDate(s string, location *time.Location) (*time.Time, error) {
	timeRaw, _, found := strings.Cut(s, "(")
	if !found {
		return nil, errors.New("could not cut string by '(' :" + timeRaw)
	}

	timeRaw = translateTimeGrEn(timeRaw)

	t, err := time.ParseInLocation("Monday 2 January 2006 - 15:04 pm", timeRaw, location)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// TranslateTimeGrEn translates Greek string time of format
// "Τετάρτη 21 Δεκεμβρίου 2022 - 11:59 μ.μ." in English
func translateTimeGrEn(dt string) string {
	parts := strings.Split(dt, " ")

	dayGR := &parts[0]
	dt = strings.Replace(dt, *dayGR, dayNamesGR[*dayGR], 1)

	monthGR := &parts[2]
	dt = strings.Replace(dt, *monthGR, monthNamesGenitiveGR[*monthGR], 1)

	periodGR := &parts[6]
	dt = strings.Replace(dt, *periodGR, periodsGR[*periodGR], 1)

	return dt
}
