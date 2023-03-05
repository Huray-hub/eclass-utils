package assignment

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/Huray-hub/eclass-utils/assignments/config"
	auth "github.com/Huray-hub/eclass-utils/authentication"
	"github.com/Huray-hub/eclass-utils/course"
)

func ExampleService_FetchAssignments_importConfigFromYaml() {
	// Import options and credentials from config.yml
	opts, creds, err := config.Import()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()
	// Not providing http.Client is fine, NewService will initialize its own
	service, err := NewService(ctx, opts, *creds, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	assignments, err := service.FetchAssignments(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Process result however you like
	_ = assignments
}

func ExampleService_FetchAssignments_yourOwnConfig() {
	// Provide your own options and credentials
	// view README.md for more info
	opts := &config.Options{
		PlainText:      false,
		IncludeExpired: false,
		ExportICS:      false,
		ExcludedAssignments: map[string][]string{
			"courseID": {
				"for people who are not registered to any lab",
				"Monday lab",
				"Monday Lab",
			},
		},
		Options: course.Options{
			ExcludedCourses: map[string]struct{}{
				"courseID":  {},
				"anotherCourseID": {},
			},
		},
	}
	creds := auth.Credentials{
		Username: "your-username",
		Password: "your-password",
	}

	// Providing your own client, authentication requires cookie Jar.
	// If you don't provide it, NewService will do it for you.
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	client := &http.Client{Jar: jar}

	ctx := context.Background()
	service, err := NewService(ctx, opts, creds, client)
	if err != nil {
		log.Fatal(err.Error())
	}
	assignments, err := service.FetchAssignments(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Process result however you like
	_ = assignments
}

func BenchmarkFetchAssignments(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opts, creds, err := config.Import()
		if err != nil {
			log.Fatal(err.Error())
		}

		opts.IncludeExpired = true

		ctx := context.Background()
		service, err := NewService(ctx, opts, *creds, nil)
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
