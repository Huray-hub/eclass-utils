package assignment

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Huray-hub/eclass-utils/assignment/config"
)

func BenchmarkFetchAssignments(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cfg, err := config.ImportDefault()
		if err != nil {
			log.Fatal(err.Error())
		}

		cfg.Options.IncludeExpired = true

		ctx := context.Background()
		service, err := NewService(ctx, *cfg, nil)
		if err != nil {
			log.Fatal(err.Error())
		}
		assignments, err := service.FetchAssignments(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%v", len(assignments))
	}
}

func TestParseNearDeadline_Tomorrow(t *testing.T) {
	t.Skip("not ready")
	// Arrange
	location, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		t.Errorf("failed to load location %v", err)
	}

	deadlineStr := "αύριο - 11:59 μ.μ.(απομένουν 1 ημέρα 3 ώρες 8 λεπτά)"
	expectedDeadline := time.Date(2022, 12, 4, 23, 59, 0, 0, location)

	// Act
	deadline, err := parseDeadline(deadlineStr, location)
	if err != nil {
		t.Errorf("failed to parse deadline: '%v'", deadline)
	}

	// Assert
	if !deadline.Equal(expectedDeadline) {
		t.Errorf("Expected: %s, Actual: %s", expectedDeadline, deadline)
	}
}

func TestParseNormalDeadline(t *testing.T) {
	t.Skip("not ready")
	// Arrange
	location, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		t.Errorf("failed to load location %v", err)
	}

	deadlineStr := "Τετάρτη 21 Δεκεμβρίου 2022 - 11:59 μ.μ.(απομένουν 19 ημέρες 3 ώρες 8 λεπτά)"
	expectedDeadline := time.Date(2022, 12, 21, 23, 59, 0, 0, location)

	// Act
	deadline, err := parseDeadline(deadlineStr, location)
	if err != nil {
		t.Errorf("failed to parse deadline: '%v'", deadline)
	}

	// Assert
	if !deadline.Equal(expectedDeadline) {
		t.Errorf("Expected: %v, Actual: %v", expectedDeadline, deadline)
	}
}
