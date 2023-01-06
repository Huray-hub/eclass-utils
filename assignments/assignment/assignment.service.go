package assignment

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Huray-hub/eclass-utils/assignments/config"
	"github.com/Huray-hub/eclass-utils/assignments/course"
	"github.com/Huray-hub/eclass-utils/assignments/login"
	"github.com/gocolly/colly"
)

var location *time.Location 

func init() {
	var err error
	location, err = time.LoadLocation("Europe/Athens")
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func Get(
	opts *config.Options,
	credentials *config.Credentials,
) ([]Assignment, error) {
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

	courses, err := course.Get(opts, c.Clone())
	if err != nil {
		return nil, err
	}

	assignments, err := getAssignments(opts, courses, c.Clone())
	if err != nil {
		return nil, err
	}

	return assignments, nil
}

func getAssignments(
	opts *config.Options, courses []course.Course, c *colly.Collector,
) ([]Assignment, error) {
	assignments := make(sortable, 0, len(courses))

	for _, crs := range courses {
		var excludedStrings []string
		if val, ok := opts.ExcludedAssignments[crs.ID]; ok {
			excludedStrings = val
		}
		apc, err := getAssignmentsPerCourse(
			opts,
			excludedStrings,
			crs,
			c.Clone(),
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, apc...)
	}

	sortAssignments(assignments)
	return assignments, nil
}

func getAssignmentsPerCourse(
	opts *config.Options,
	excludedStrings []string,
	course course.Course,
	c *colly.Collector,
) ([]Assignment, error) {
	assignments := make([]Assignment, 0, 10)

	isExcluded := func(assignment *Assignment) bool {
		if !opts.IncludeExpired && assignment.Deadline.Before(time.Now().In(location)) {
			return true
		}

		if excludedStrings == nil {
			return false
		}

		for _, v := range excludedStrings {
			if strings.Contains(assignment.Title, v) {
				return true
			}
		}

		return false
	}

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL,
			"failed with response:", r, "\nError:", err)
		log.Fatal(err.Error())
	})

	c.OnHTML(
		"table#assignment_table tbody tr",
		func(h *colly.HTMLElement) {
			if h == nil {
				return
			}

			tds := make([]*colly.HTMLElement, 0, 4)
			h.ForEach("td", func(_ int, h2 *colly.HTMLElement) {
				tds = append(tds, h2)
			})

			assignment, err := newAssignment(tds, &course, location)
			if err != nil {
				return
			}

			if isExcluded(assignment) {
				return
			}

			assignments = append(assignments, *assignment)
		},
	)

	finalURL, err := course.PrepareAssignmentsURL(opts.BaseDomain)
	if err != nil {
		return nil, err
	}

	err = c.Visit("https://" + finalURL)
	if err != nil {
		return nil, err
	}

	return assignments, nil
}
