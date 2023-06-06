package assignment

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Huray-hub/eclass-utils/assignment/config"
	"github.com/Huray-hub/eclass-utils/course"
	"github.com/PuerkitoBio/goquery"
	"errors"
)

type Assignment struct {
	ID       string
	Course   *course.Course
	Title    string
	Deadline *time.Time
	IsSent   bool
}

func (a *Assignment) String() string {
	var deadlineMsg string
	if a.Deadline == nil {
		deadlineMsg = NoDeadline
	} else {
		deadlineMsg =  a.Deadline.Format("02/01/2006 15:04")
	}

	return fmt.Sprintf(
		"%v,%v,%v,%v,%v,%v",
		a.Course.ID,
		a.Course.Name,
		a.ID,
		a.Title,
		deadlineMsg,
		a.IsSent,
	)
}

func newAssignment(
	tds []*goquery.Selection,
	course *course.Course,
	location *time.Location,
) (Assignment, error) {
	deadline, err := parseDeadline(tds[1].Text(), location)
	if err != nil {
		return Assignment{}, err
	}
	id, err := parseID(tds[0])
	if err != nil {
		return Assignment{}, err
	}

	return Assignment{
		ID:       id,
		Course:   course,
		Title:    strings.TrimSpace(tds[0].Text()),
		Deadline: deadline,
		IsSent:   parseIsSent(tds[2]),
	}, nil
}

func parseID(td *goquery.Selection) (string, error) {
	uri, ok := td.Find("a").Attr("href")
	if !ok {
		return "", errors.New("could not parse ID from 'a href'")
	}

	uriValues, err := url.ParseQuery(uri)
	if err != nil {
		return "", err
	}

	id := uriValues.Get("id")
	if _, err := strconv.Atoi(id); err != nil {
		return "", fmt.Errorf("ID: %v is not a valid string", id)
	}

	return id, nil
}

func parseDeadline(dl string, location *time.Location) (*time.Time, error) {
	if dl == NoDeadline {
		return nil, nil
	}

	t, err := parseTime(dl, location)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func parseIsSent(s *goquery.Selection) bool {
	return s.Children().First().HasClass("fa-check-square-o")
}

// SortByDeadline function sorts assignments by descending deadline
func SortByDeadline(a []Assignment) {
	sort.Slice(a, func(i, j int) bool {
		if a[i].Deadline == nil {
			return false
		}
		if a[j].Deadline == nil {
			return true
		}
		return a[i].Deadline.Before(*a[j].Deadline)
	})
}

// PrepareURL method prepares URL for assignments' own page
func (a *Assignment) PrepareURL(baseURL string) (string, error) {
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

// IsExcluded method determines if the assignment should be excluded from final result
func (a Assignment) IsExcluded(
	opts config.Options,
	courseID string,
	location *time.Location,
) bool {
	if !opts.IncludeExpired && a.Deadline != nil && a.Deadline.Before(time.Now().In(location)) {
		return true
	}

	for _, excludedString := range opts.ExcludedAssignments[courseID] {
		if strings.Contains(a.Title, excludedString) {
			return true
		}
	}

	return false
}
