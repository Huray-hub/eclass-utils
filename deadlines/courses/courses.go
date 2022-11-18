package courses

import (
	"net/url"
	"strings"

	"github.com/Huray-hub/eclass-utils/deadlines/config"
	"github.com/gocolly/colly"
)

type Course struct {
	ID   string
	Name string
	URL  string
}

func newCourse(name, url string) Course {
	return Course{
		ID:   extractID(url),
		Name: name,
		URL:  url,
	}
}

func (c Course) PrepareAssignmentsURL(baseURL string) (string, error) {
	finalURL, err := url.Parse(baseURL + "/modules/work/")
	if err != nil {
		return "", err
	}

	values := finalURL.Query()
	values.Add("course", c.ID)
	finalURL.RawQuery = values.Encode()

	return finalURL.String(), nil
}

func extractID(url string) string {
	urlParts := strings.Split(url, "/")
	return urlParts[len(urlParts)-2]
}

func Get(
	opts *config.Options,
	c *colly.Collector,
) ([]Course, error) {
	courses := make([]Course, 0, 10)

	isExcluded := func(course Course) bool {
		if len(opts.ExcludedCourses) == 0 {
			return false
		}
		for _, id := range opts.ExcludedCourses {
			if id == course.ID {
				return true
			}
		}
		return false
	}

	c.OnHTML("#main-content table.table-default tbody tr a",
		func(h *colly.HTMLElement) {
			if len(h.Text) > 0 {
				course := newCourse(h.Text, h.Attr("href"))

				if isExcluded(course) {
					return
				}

				courses = append(courses, newCourse(h.Text, h.Attr("href")))
			}
		})

	err := c.Visit("https://" + opts.BaseDomain + "/main/my_courses.php")
	if err != nil {
		return nil, err
	}

	return courses, nil
}
