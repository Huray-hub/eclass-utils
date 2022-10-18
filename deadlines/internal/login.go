package internal

import (
	"fmt"

	"github.com/gocolly/colly"
)

func headHomepage(url string, c *colly.Collector) {
	c.Visit("https://" + url)
}

func Login(url string, creds Creds, c *colly.Collector) error {
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	headHomepage(url, c)
	postLogin(url, creds, c)

	return nil
}

func postLogin(url string, creds Creds, c *colly.Collector) {
	body := make(map[string]string, 3)

	body["uname"] = creds.Username
	body["pass"] = creds.Password
	body["submit"] = ""

	c.Post("https://"+url, body)
}
