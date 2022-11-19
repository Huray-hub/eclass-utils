package course

import (
	"net/url"
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
