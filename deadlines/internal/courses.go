package internal

import (
	"strings"

	"github.com/gocolly/colly"
)

type course struct {
	Code string
	Name string
	Url  string
}

func newCourse(name, url string) course {
	return course{
		Code: extractCourseCode(url),
		Name: name,
		Url:  url,
	}
}

func extractCourseCode(url string) string {
	urlParts := strings.Split(url, "/")
	return urlParts[len(urlParts)-2]
}

func GetEnrolledCourses(c *colly.Collector) []course {
	courses := make([]course, 0, 10)

	c.OnHTML("#main-content table.table-default tbody tr a",
		func(h *colly.HTMLElement) {
			if len(h.Text) > 0 {
				courses = append(courses, newCourse(h.Text, h.Attr("href")))
			}
		})

	c.Visit(BASE_URL + "/main/my_courses.php")

	return courses
}
