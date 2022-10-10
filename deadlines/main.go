package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

const (
	BASE_DOMAIN string = "eclass.uniwa.gr"
	BASE_URL    string = "https://eclass.uniwa.gr"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains(BASE_DOMAIN),
	)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL,
			"failed with response:", r, "\nError:", err)
	})

	Login(c)
	courses := GetEnrolledCourses(c.Clone())
	assignments := FetchAssignments(c.Clone(), courses)

	fmt.Println(assignments)
}
