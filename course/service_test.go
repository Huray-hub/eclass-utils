package course_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Huray-hub/eclass-utils/assignments/config"
	auth "github.com/Huray-hub/eclass-utils/authentication"
	"github.com/Huray-hub/eclass-utils/course"
)

// TODO: replace assignments dependency
func BenchmarkGetEnrolledCourses(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opts, creds, err := config.Import()
		if err != nil {
			log.Fatal(err.Error())
		}

		ctx := context.Background()
		client, err := auth.Login(ctx, "https://"+opts.BaseDomain, *creds, nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		courses, err := course.GetEnrolled(ctx, opts.Options, client)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%v", len(courses))
	}
}
