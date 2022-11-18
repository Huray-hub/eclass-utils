package assignments

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Huray-hub/eclass-utils/deadlines/config"
	"github.com/Huray-hub/eclass-utils/deadlines/courses"
	"github.com/gocolly/colly"
)

type Assignment struct {
	ID       string
	Course   *courses.Course
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
	course *courses.Course,
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

func Get(
	opts *config.Options, courses []courses.Course, c *colly.Collector,
) ([]Assignment, error) {
	assignments := make(sortableSlice, 0, len(courses))

	for _, course := range courses {
		var filteredOutKeywords []string
		if val, ok := opts.ExcludedAssignmentsByKeyword[course.ID]; ok {
			filteredOutKeywords = val
		}
		apc, err := fetchAssignmentsPerCourse(
			opts,
			filteredOutKeywords,
			course,
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

func sortAssignments(a sortableSlice) {
	sort.Sort(a)
}

func fetchAssignmentsPerCourse(
	opts *config.Options,
	filteredOutKeywords []string,
	course courses.Course,
	c *colly.Collector,
) ([]Assignment, error) {
	assignments := make([]Assignment, 0, 10)

	location, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		return nil, err
	}

	isExcluded := func(assignment *Assignment) bool {
		if opts.IgnoreExpired && assignment.Deadline.Before(time.Now().In(location)) {
			return true
		}

		if filteredOutKeywords == nil {
			return false
		}

		for _, v := range filteredOutKeywords {
			if strings.Contains(assignment.Title, v) {
				return true
			}
		}

		return false
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

func (a Assignment) PrepareURL(
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
