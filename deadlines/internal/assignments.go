package internal

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Assignment struct {
	ID       string
	Course   *Course
	Title    string
	Deadline time.Time
	IsSent   bool
}

func (a *Assignment) String() string {
	return fmt.Sprintf(
		"%v,%v,%v,%v",
		a.Course.Name,
		a.Title,
		a.Deadline.String(),
		a.IsSent,
	)
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

func newAssignment(
	tds []*colly.HTMLElement,
	course *Course,
	location *time.Location,
) (*Assignment, error) {
	deadline, err := parseDeadline(tds[1].Text, location)
	if err != nil {
		return nil, err
	}
	id, err := parseID(tds[0])
	if err != nil {
		return nil, err
	}

	return &Assignment{
		ID:       id,
		Course:   course,
		Title:    tds[0].Text,
		Deadline: deadline,
		IsSent:   parseIsSent(tds[2]),
	}, nil
}

func parseID(td *colly.HTMLElement) (string, error) {
	uri := td.ChildAttr("a", "href")
	if uri == "" {
		return "", fmt.Errorf("could not parse assignment's ID from url: %v", uri)
	}

	id := strings.Split(uri, "id=")[1]

	if _, err := strconv.Atoi(id); err != nil {
		return "", fmt.Errorf("ID: %v is not a valid string", uri)
	}

	return id, nil
}

func parseDeadline(dl string, location *time.Location) (time.Time, error) {
	// deadline, _ := time.ParseInLocation("02-01-2006 15:04:05", "30-11-2022 23:55:00", location)
	dt, err := time.ParseInLocation(
		"02-01-2006 15:04:05",
		strings.Split(dl, "(")[0],
		location,
	)
	if err != nil {
		return time.Time{}, err
	}
	return dt, nil
}

func parseIsSent(h *colly.HTMLElement) bool {
	return h.DOM.Children().First().HasClass("fa-check-square-o")
}

func FetchAssignments(
	url string, courses []Course, c *colly.Collector,
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
	course Course,
	c *colly.Collector,
) ([]Assignment, error) {
	assignments := make([]Assignment, 0, 10)

	location, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		return nil, err
	}

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

			assignment, err2 := newAssignment(tds, &course, location)
			if err2 != nil {
				return
			}

			assignments = append(assignments, *assignment)
		},
	)

	finalURL, err := prepareCourseURL(url, course)
	if err != nil {
		return nil, err
	}

	err = c.Visit("https://" + finalURL)
	if err != nil {
		return nil, err
	}

	return assignments, nil
}

func (a Assignment) prepareAssignmentURL(
	baseURL string,
) (string, error) {
	finalURL, err := url.Parse(baseURL + "/modules/work/index.php")
	if err != nil {
		return "", err
	}

	values := finalURL.Query()
	values.Add("course", a.Course.ID)
	values.Add("id", a.ID)
	finalURL.RawQuery = values.Encode()

	return finalURL.String(), nil
}

func prepareCourseURL(baseURL string, course Course) (string, error) {
	finalURL, err := url.Parse(baseURL + "/modules/work/")
	if err != nil {
		return "", err
	}

	values := finalURL.Query()
	values.Add("course", course.ID)
	finalURL.RawQuery = values.Encode()

	return finalURL.String(), nil
}

func FilterExpiredDeadlines(assignments []Assignment) ([]Assignment, error) {
	location, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		return nil, err
	}
	timeNow := time.Now().In(location)

	filteredAssignments := make([]Assignment, 0, len(assignments))
	for _, assignment := range assignments {
		if assignment.Deadline.After(timeNow) {
			filteredAssignments = append(filteredAssignments, assignment)
		}
	}
	return filteredAssignments, nil
}
