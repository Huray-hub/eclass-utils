package calendar_test

import (
	"testing"
	"time"

	"github.com/Huray-hub/eclass-utils/assignments/assignment"
	"github.com/Huray-hub/eclass-utils/assignments/calendar"
	"github.com/Huray-hub/eclass-utils/course"
)

func TestExport(t *testing.T) {
	t.Skip("currently I use this only as a shorcut to my workflow")
	// Arrange
	baseDomain := "eclass.uniwa.gr"
	course := &course.Course{ID: "ICE262", Name: "ΑΝΑΚΤΗΣΗ ΠΛΗΡΟΦΟΡΙΑΣ"}
	location, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		t.Fatal(err.Error())
	}

	assignments := [2]assignment.Assignment{
		{
			ID:     "24692",
			Course: course,
			Deadline: func(location *time.Location) time.Time {
				var deadline time.Time
				deadline, err = time.ParseInLocation(
					"02-01-2006 15:04:05",
					"30-11-2022 23:55:00",
					location,
				)
				if err != nil {
					t.Error("cannot parse string to local deadline")
				}
				return deadline
			}(location),
			IsSent: false,
			Title:  "Άσκηση 1 (τμήματα Τετάρτης)",
		},
		{
			ID:     "15207",
			Course: course,
			Deadline: func(location *time.Location) time.Time {
				var deadline time.Time
				deadline, err = time.ParseInLocation(
					"02-01-2006 15:04:05",
					"28-11-2022 23:55:00",
					location,
				)
				if err != nil {
					t.Error("cannot parse string to local deadline")
				}
				return deadline
			}(location),
			IsSent: false,
			Title:  "Άσκηση 1 (τμήματα Δευτέρας)",
		},
	}

	// Act
	res, err := calendar.Export(assignments[:], baseDomain)

	// Assert
	if err != nil {
		t.Errorf(err.Error())
	}

	if res == "" {
		t.Errorf("Empty result\n")
	}
}
