package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func headHomepage(c *colly.Collector) {
	c.Visit(BASE_URL)
}

func Login(c *colly.Collector) error {
	cfg, err := GetConfiguration()
	if err != nil {
		return err
	}

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	headHomepage(c)
	postLogin(c, cfg)

	return nil
}

func postLogin(c *colly.Collector, cfg map[string]string) {
	cfg["submit"] = ""

	c.Post(BASE_URL, cfg)
}
