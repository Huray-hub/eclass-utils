package internal

import (
	"strings"

	"github.com/gocolly/colly"
)

type Course struct {
	ID   string
	Name string
	URL  string
}

func newCourse(name, url string) Course {
	return Course{
		ID:   extractCourseID(url),
		Name: name,
		URL:  url,
	}
}

func extractCourseID(url string) string {
	urlParts := strings.Split(url, "/")
	return urlParts[len(urlParts)-2]
}

func GetEnrolledCourses(url string, c *colly.Collector) ([]Course, error) {
	courses := make([]Course, 0, 10)

	c.OnHTML("#main-content table.table-default tbody tr a",
		func(h *colly.HTMLElement) {
			if len(h.Text) > 0 {
				courses = append(courses, newCourse(h.Text, h.Attr("href")))
			}
		})

	err := c.Visit("https://" + url + "/main/my_courses.php")
	if err != nil {
		return nil, err
	}

	return courses, nil
}
