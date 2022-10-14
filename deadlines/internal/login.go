package internal

import (
	"fmt"

	"github.com/gocolly/colly"
)

func headHomepage(c *colly.Collector) {
	c.Visit(BASE_URL)
}

func Login(baseUrl string,cfg map[string]string, c *colly.Collector) error {
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	headHomepage(c)
	postLogin(c, cfg)

	return nil
}

func postLogin(c *colly.Collector, cfg map[string]string) {
    body := make(map[string]string, 3)

    body["uname"] = cfg["uname"]
    body["pass"] = cfg["pass"]
	body["submit"] = ""

	c.Post(BASE_URL, body)
}
