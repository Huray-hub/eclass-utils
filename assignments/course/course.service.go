package course

import (
	"github.com/Huray-hub/eclass-utils/assignments/config"
	"github.com/gocolly/colly"
)

func Get(
	opts *config.Options,
	c *colly.Collector,
) ([]Course, error) {
	courses := make([]Course, 0, 10)

	isExcluded := func(course Course) bool {
		if _, ok := opts.ExcludedCourses[course.ID]; ok {
			return true
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
