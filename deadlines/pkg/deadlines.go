package deadlines

import (
	"fmt"

	in "github.com/Huray-hub/eclass-utils/deadlines/internal"
	"github.com/gocolly/colly"
)

func Deadlines(opts *in.Options, creds *in.Creds) ([]in.Assignment, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(opts.BaseDomain),
	)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL,
			"failed with response:", r, "\nError:", err)
	})

	err := in.Login(opts.BaseDomain, *creds, c)
	if err != nil {
		return nil, err
	}

	courses := in.GetEnrolledCourses(opts.BaseDomain, c.Clone())

	assignments, err := in.FetchAssignments(opts.BaseDomain, *courses,c.Clone())
	if err != nil {
		return nil, err
	}

	return assignments, nil
}
