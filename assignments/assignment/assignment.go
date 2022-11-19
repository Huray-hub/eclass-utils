package assignment

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Huray-hub/eclass-utils/assignments/course"
	"github.com/gocolly/colly"
)

type Assignment struct {
	ID       string
	Course   *course.Course
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

type sortable []Assignment

func (a sortable) Len() int {
	return len(a)
}

func (a sortable) Less(i, j int) bool {
	return a[i].Deadline.Before(a[j].Deadline)
}

func (a sortable) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func newAssignment(
	tds []*colly.HTMLElement,
	course *course.Course,
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

func sortAssignments(a sortable) {
	sort.Sort(a)
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
