package internal

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Assignment struct {
	course   string
	title    string
	deadline time.Time
	isSent   bool
}

func (a *Assignment) String() string {
	return fmt.Sprintf("%v %v %v %v", a.course, a.title, a.deadline, a.isSent)
}

func newAssignment(tds []*colly.HTMLElement, courseName string) (Assignment, error) {
	dl, err := parseDeadline(tds[1].Text)
	if err != nil {
		return Assignment{}, err
	}

	return Assignment{
		course:   courseName,
		title:    tds[0].Text,
		deadline: dl,
		isSent:   parseIsSent(tds[2]),
	}, nil
}

func parseDeadline(dl string) (time.Time, error) {
	dt, err := time.Parse("02-01-2006 15:04:05", strings.Split(dl, "(")[0])
	if err != nil {
		return time.Time{}, err
	}
	return dt, nil
}

func parseIsSent(h *colly.HTMLElement) bool {
	return h.DOM.Children().First().HasClass("fa-check-square-o")
}

func FetchAssignments(c *colly.Collector, courses []course) ([]Assignment, error) {
	assignments := make([]Assignment, 0, len(courses))

	for _, course := range courses {
		apc, err := fetchAssignmentsPerCourse(c.Clone(), course)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, apc...)
	}

	return assignments, nil
}

func fetchAssignmentsPerCourse(c *colly.Collector, course course) ([]Assignment, error) {
	assignments := make([]Assignment, 0, 10)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL,
			"failed with response:", r, "\nError:", err)
		// log.Fatal(err.Error())
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

			newAss, err := newAssignment(tds, course.Name)
			if err != nil {
				return
			}

			assignments = append(assignments, newAss)
		})

	finalUrl, err := prepareCourseUrl(course)
	if err != nil {
		return nil, err
	}

	c.Visit(finalUrl)

	return assignments, nil
}

func prepareCourseUrl(course course) (string, error) {
	url, err := url.Parse("" + "/modules/work/")
	if err != nil {
		return "", err
	}

	values := url.Query()
	values.Add("course", course.Code)
	url.RawQuery = values.Encode()

	return url.String(), nil
}
