package deadlines

import (
	"fmt"

	"github.com/Huray-hub/eclass-utils/deadlines/assignments"
	"github.com/Huray-hub/eclass-utils/deadlines/config"
	"github.com/Huray-hub/eclass-utils/deadlines/courses"
	"github.com/Huray-hub/eclass-utils/deadlines/login"
	"github.com/gocolly/colly"
)

func Get(
	opts *config.Options,
	credentials *config.Credentials,
) ([]assignments.Assignment, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(opts.BaseDomain),
	)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL,
			"failed with response:", r, "\nError:", err)
	})

	err := login.Login(opts.BaseDomain, *credentials, c)
	if err != nil {
		return nil, err
	}

	courses, err := courses.Get(opts, c.Clone())
	if err != nil {
		return nil, err
	}

	assignments, err := assignments.Get(opts, courses, c.Clone())
	if err != nil {
		return nil, err
	}

	return assignments, nil
}
