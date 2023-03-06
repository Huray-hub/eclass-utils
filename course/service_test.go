package course_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"testing"

	"github.com/Huray-hub/eclass-utils/auth"
	"github.com/Huray-hub/eclass-utils/course"

	"github.com/Huray-hub/eclass-utils/assignment/config"
)

func ExampleGetEnrolled() {
	opts := course.Options{
		BaseDomain: "eclass.your-university.gr",
		ExcludedCourses: map[string]struct{}{
			"courseID":        {},
			"anotherCourseID": {},
		},
	}
	creds := auth.Credentials{
		Username: "your-username",
		Password: "your-password",
	}

	// Providing an http.client
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	client := &http.Client{Jar: jar}
	// Login
	client, err = auth.Login(context.Background(), "https://"+opts.BaseDomain, creds, client)
	if err != nil {
		return
	}

	courses, err := course.GetEnrolled(context.Background(), opts, client)
	if err != nil {
		return
	}

	// Process result however you like
	_ = courses
}

// TODO: replace assignments dependency
func BenchmarkGetEnrolledCourses(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cfg, err := config.ImportDefault()
		if err != nil {
			log.Fatal(err.Error())
		}

		ctx := context.Background()
		client, err := auth.Login(ctx, "https://"+cfg.Options.BaseDomain, cfg.Credentials, nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		courses, err := course.GetEnrolled(ctx, cfg.Options.Options, client)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%v", len(courses))
	}
}
