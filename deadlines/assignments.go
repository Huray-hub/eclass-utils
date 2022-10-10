package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type assignment struct {
	course   string
	title    string
	deadline time.Time
	isSent   bool
}

func (a *assignment) String() string {
	return fmt.Sprintf("%v %v %v %v", a.course, a.title, a.deadline, a.isSent)
}

func newAssignment(tds []*colly.HTMLElement, courseName string) assignment {
	return assignment{
		course:   courseName,
		title:    tds[0].Text,
		deadline: parseDeadline(tds[1].Text),
		isSent:   parseIsSent(tds[2]),
	}
}

func parseDeadline(dl string) time.Time {
	dt, err := time.Parse("02-01-2006 15:04:05", strings.Split(dl, "(")[0])
	if err != nil {
		// TODO: handle accordingly
		// return time.Now(), err
	}
	return dt
}

func parseIsSent(h *colly.HTMLElement) bool {
	return h.DOM.Children().First().HasClass("fa-check-square-o")
}

func FetchAssignments(c *colly.Collector, courses []course) []assignment {
	assignments := make([]assignment, 0, len(courses))

	for _, course := range courses {
		assignments = append(assignments, fetchAssignmentsPerCourse(c.Clone(), course)...)
	}

	return assignments
}

func fetchAssignmentsPerCourse(c *colly.Collector, course course) []assignment {
	assignments := make([]assignment, 0, 10)

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

			assignments = append(assignments, newAssignment(tds, course.Name))
		})

	finalUrl, err := prepareCourseUrl(course)
	if err != nil {
		return nil
	}

	c.Visit(finalUrl)

	return assignments
}

func prepareCourseUrl(course course) (string, error) {
	url, err := url.Parse(BASE_URL + "/modules/work/")
	if err != nil {
		return "", err
	}

	values := url.Query()
	values.Add("course", course.Code)
	url.RawQuery = values.Encode()

	return url.String(), nil
}
