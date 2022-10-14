package deadlines

import (
	"fmt"
	"log"

	"github.com/Huray-hub/eclass-utils/deadlines/internal"
	"github.com/gocolly/colly"
)

const (
	BASE_DOMAIN string = "eclass.uniwa.gr"
	BASE_URL    string = "https://eclass.uniwa.gr"
)

func Deadlines() {
	cfg, err := internal.GetConfiguration()

	baseDomain := cfg["domain"]

	if err != nil {
		log.Fatal(err.Error())
	}

	c := colly.NewCollector(
		colly.AllowedDomains(baseDomain),
	)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL,
			"failed with response:", r, "\nError:", err)
	})

	err := internal.Login(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	courses := internal.GetEnrolledCourses(c.Clone())

	assignments, err := internal.FetchAssignments(c.Clone(), courses)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(assignments)
}
