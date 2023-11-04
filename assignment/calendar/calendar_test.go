package calendar_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Huray-hub/eclass-utils/assignment"
	"github.com/Huray-hub/eclass-utils/assignment/calendar"
	"github.com/Huray-hub/eclass-utils/course"
)

func TestExport(t *testing.T) {
	// Arrange
	baseDomain := "eclass.uniwa.gr"
	course := &course.Course{ID: "ICE262", Name: "ΑΝΑΚΤΗΣΗ ΠΛΗΡΟΦΟΡΙΑΣ"}
	location, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		t.Fatalf("arrange phase - %s", err)
	}

	assignments := [2]assignment.Assignment{
		{
			ID:     "24692",
			Course: course,
			Deadline: func(location *time.Location) *time.Time {
				var deadline time.Time
				deadline, err = time.ParseInLocation(
					"02-01-2006 15:04:05",
					"30-11-2022 23:55:00",
					location,
				)
				if err != nil {
					t.Fatalf("arrange phase - cannot parse string to local deadline: %s", err)
				}
				return &deadline
			}(location),
			IsSent: false,
			Title:  "Άσκηση 1 (τμήματα Τετάρτης)",
		},
		{
			ID:     "15207",
			Course: course,
			Deadline: func(location *time.Location) *time.Time {
				var deadline time.Time
				deadline, err = time.ParseInLocation(
					"02-01-2006 15:04:05",
					"28-11-2022 23:55:00",
					location,
				)
				if err != nil {
					t.Fatalf("arrange phase - cannot parse string to local deadline: %s", err)
				}
				return &deadline
			}(location),
			IsSent: false,
			Title:  "Άσκηση 1 (τμήματα Δευτέρας)",
		},
	}

	tempDir := t.TempDir()

	expectedICSFileName := fmt.Sprintf("assignments_%s.ics", time.Now().Format("01-02-2006"))

	// Act
	filePath, err := calendar.Export(assignments[:], baseDomain, tempDir)
	// Assert
	if err != nil {
		t.Fatal(err)
	}

	if filePath == "" {
		t.Fatal("empty filePath")
	}

	actualICSFileName, _ := strings.CutPrefix(filePath, tempDir)
	if expectedICSFileName == actualICSFileName {
		t.Fatalf("expected: %s, actual: %s", expectedICSFileName, actualICSFileName)
	}
}
