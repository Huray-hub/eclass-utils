package calendar_test

import (
	"testing"
	"time"

	in "github.com/Huray-hub/eclass-utils/deadlines/internal"
)

func TestExportICS(t *testing.T) {
	// Arrange
	baseDomain := "eclass.uniwa.gr"
	course := &in.Course{ID: "ICE262", Name: "ΑΝΑΚΤΗΣΗ ΠΛΗΡΟΦΟΡΙΑΣ"}
	location, _ := time.LoadLocation("Europe/Athens")

	assignments := [2]in.Assignment{
		{
			ID:     "24692",
			Course: course,
			Deadline: func(location *time.Location) time.Time {
				deadline, err := time.ParseInLocation(
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
				deadline, err := time.ParseInLocation(
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
	res, err := in.ExportICS(assignments[:], baseDomain)

	// Assert
	if err != nil {
		t.Errorf(err.Error())
	}

	if res == "" {
		t.Errorf("Empty result\n")
	}
}
