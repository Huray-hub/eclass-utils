package main

import (
	"time"

	"github.com/gocolly/colly"
)

type deadline struct {
	course string
	title  string
	date   time.Time
	isSent bool
}

func newDeadline() deadline {
	return deadline{}
}

func FetchAssignments(c *colly.Collector, courses []course) []deadline {
	dl := make([]deadline, 0, len(courses)/2*3)

	for _, course := range courses {
		dl = append(dl, fetchAssignmentsPerCourse(c.Clone(), course))
	}

	return dl
}

func fetchAssignmentsPerCourse(c *colly.Collector, course course) deadline {
	c.OnHTML(
		"table#assignment_table tbody tr a",
		func(h *colly.HTMLElement) {
			if h == nil {
				return
			}

			// if len(h.Text) > 0 {
			// 	courses = append(courses, newCourse(h.Text, h.Attr("href")))
			// }
		})

	c.Visit(BASE_URL + "/modules/work/")

	return deadline{}
}
