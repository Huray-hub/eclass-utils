package course_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Huray-hub/eclass-utils/auth"
	"github.com/Huray-hub/eclass-utils/course"

	"github.com/Huray-hub/eclass-utils/assignment/config"
)

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
