package internal

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Assignment struct {
	Course   string
	Title    string
	Deadline time.Time
	IsSent   bool
}

func (a *Assignment) String() string {
	return fmt.Sprintf("%v,%v,%v,%v", a.Course, a.Title, a.Deadline.String(), a.IsSent)
}

type sortableSlice []Assignment

func (p sortableSlice) Len() int {
	return len(p)
}

func (p sortableSlice) Less(i, j int) bool {
	return p[i].Deadline.Before(p[j].Deadline)
}

func (p sortableSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func newAssignment(tds []*colly.HTMLElement, course string) (*Assignment, error) {
	dl, err := parseDeadline(tds[1].Text)
	if err != nil {
		return nil, err
	}

	return &Assignment{
		Course:   course,
		Title:    tds[0].Text,
		Deadline: dl,
		IsSent:   parseIsSent(tds[2]),
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

func FetchAssignments(
	url string, courses []course, c *colly.Collector,
) ([]Assignment, error) {
	assignments := make(sortableSlice, 0, len(courses))

	for _, course := range courses {
		apc, err := fetchAssignmentsPerCourse(url, course, c.Clone())
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, apc...)
	}

	sortAssignments(assignments)
	return assignments, nil
}

func sortAssignments(a sortableSlice) {
	sort.Sort(a)
}

func fetchAssignmentsPerCourse(
	url string,
	course course,
	c *colly.Collector,
) ([]Assignment, error) {
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

			assignments = append(assignments, *newAss)
		})

	finalUrl, err := prepareCourseUrl(url, course)
	if err != nil {
		return nil, err
	}

	err = c.Visit("https://" + finalUrl)
	if err != nil {
		return nil, err
	}

	return assignments, nil
}

func prepareCourseUrl(baseUrl string, course course) (string, error) {
	url, err := url.Parse(baseUrl + "/modules/work/")
	if err != nil {
		return "", err
	}

	values := url.Query()
	values.Add("course", course.Code)
	url.RawQuery = values.Encode()

	return url.String(), nil
}
