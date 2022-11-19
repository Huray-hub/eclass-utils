package login

import (
	"fmt"

	"github.com/Huray-hub/eclass-utils/assignments/config"
	"github.com/gocolly/colly"
)

func headHomepage(url string, c *colly.Collector) error {
	err := c.Visit("https://" + url)
	if err != nil {
		return err
	}
	return nil
}

func Login(url string, credentials config.Credentials, c *colly.Collector) error {
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(
			"Request URL:",
			r.Request.URL,
			"failed with response:",
			r,
			"\nError:",
			err,
		)
	})

	err := headHomepage(url, c)
	if err != nil {
		return err
	}

	err = postLogin(url, credentials, c)
	if err != nil {
		return err
	}

	return nil
}

func postLogin(url string, credentials config.Credentials, c *colly.Collector) error {
	body := make(map[string]string, 3)

	body["uname"] = credentials.Username
	body["pass"] = credentials.Password
	body["submit"] = ""

	err := c.Post("https://"+url, body)
	if err != nil {
		return err
	}

	return nil
}
