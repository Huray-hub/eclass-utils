package course

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

type Course struct {
	ID   string
	Name string
	URL  string
}

func newCourse(name, url string) Course {
	return Course{
		ID:   extractID(url),
		Name: strings.TrimSpace(name),
		URL:  url,
	}
}

func (crs Course) String() string {
	return fmt.Sprintf("%v,%v,%v", crs.ID, crs.Name, crs.URL)
}

func (crs Course) IsExcluded(opts Options) bool {
	if _, ok := opts.ExcludedCourses[crs.ID]; ok {
		return true
	}
	return false
}

func (crs Course) PrepareAssignmentsURL(baseURL string) (string, error) {
	finalURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	finalURL = finalURL.JoinPath("modules", "work")

	values := finalURL.Query()
	values.Add("course", crs.ID)
	finalURL.RawQuery = values.Encode()

	return finalURL.String(), nil
}

func extractID(url string) string {
	return path.Base(url)
}
