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

func (c *Course) String() string {
	return fmt.Sprintf("%v,%v,%v", c.ID, c.Name, c.URL)
}

func (c *Course) PrepareAssignmentsURL(baseURL string) (string, error) {
	finalURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	finalURL = finalURL.JoinPath("modules", "work")

	values := finalURL.Query()
	values.Add("course", c.ID)
	finalURL.RawQuery = values.Encode()

	return finalURL.String(), nil
}

func extractID(url string) string {
	return path.Base(url)
}
