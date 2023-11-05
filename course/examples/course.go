package examples

import (
	"context"
	"net/http"
	"net/http/cookiejar"

	"github.com/Huray-hub/eclass-utils/auth"
	"github.com/Huray-hub/eclass-utils/course"
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
