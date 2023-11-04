package course

import (
	"context"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Options struct {
	BaseDomain          string              `yaml:"baseDomain"`
	OnlyFavoriteCourses bool                `yaml:"onlyFavoriteCourses"`
	ExcludedCourses     map[string]struct{} `yaml:"excludedCourses"`
}

// TODO: this implementation is not autonomous and needs revisit
func GetEnrolled(ctx context.Context, opts Options, client *http.Client) ([]Course, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://"+opts.BaseDomain+"/main/my_courses.php",
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			panic("could not close response body")
		}
	}()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// The max number of courses per semester is 7. Till half the semester students expect
	// courses' grades from previous semester's exam, so they are enrolled to maximum 14
	// courses in total
	courses := make([]Course, 0, 14)

	doc.Find("#main-content table.table-default tbody tr").
		Each(func(_ int, s *goquery.Selection) {
			firsttd := s.Find("td:first-child a")

			name := firsttd.Text()
			href, _ := firsttd.Attr("href")
			isFavorite := s.Find("td:nth-child(2) a:first-child span").HasClass("fa-star")

			course := newCourse(name, href, isFavorite)
			if course.IsExcluded(opts) {
				return
			}
			courses = append(courses, course)
		})

	return courses, nil
}
