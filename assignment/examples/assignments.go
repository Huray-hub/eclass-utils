package examples

import (
	"context"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/Huray-hub/eclass-utils/assignment"
	"github.com/Huray-hub/eclass-utils/assignment/config"
	"github.com/Huray-hub/eclass-utils/auth"
	"github.com/Huray-hub/eclass-utils/course"
)

// TODO: organize examples in separate packages
func FetchAssignmentsImportConfigFromYaml() {
	// Import options and credentials from config.yml
	cfg, err := config.ImportDefault()
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()
	// Not providing http.Client is fine, NewService will initialize its own
	service, err := assignment.NewService(ctx, *cfg, nil)
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

func FetchAssignmentsYourOwnConfig() {
	// Provide your own options and credentials
	// view README.md for more info
	cfg := config.Config{
		Credentials: auth.Credentials{
			Username: "your-username",
			Password: "your-password",
		},
		Options: config.Options{
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
					"courseID":        {},
					"anotherCourseID": {},
				},
			},
		},
	}

	// Providing your own client, authentication requires cookie Jar.
	// If you don't provide it, NewService will do it for you.
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	client := &http.Client{Jar: jar}

	ctx := context.Background()
	service, err := assignment.NewService(ctx, cfg, client)
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
