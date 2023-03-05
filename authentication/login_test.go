package authentication_test

import (
	"context"
	"net/http"
	"net/http/cookiejar"

	auth "github.com/Huray-hub/eclass-utils/authentication"
)

func ExampleLogin_withClient() {
	// Your credentials
	creds := auth.Credentials{
		Username: "your eclass username",
		Password: "your eclass password",
	}

	// Instantiate an empty cookie jar and pass to http.Client
	// If you don't provide one, Login function will.
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	client := &http.Client{
		Jar: jar,
	}

	// Your university domainURL
	domainURL := "https://eclass.youruniversity.gr"

	// Create session
	_, err = auth.Login(context.Background(), domainURL, creds, client)
	if err != nil {
		return
	}

	// Do something as authenticated user
	_, err = client.Get(domainURL + "/main/my_courses.php")
	if err != nil {
		return
	}
}

func ExampleLogin_noClient() {
	// Your credentials
	creds := auth.Credentials{
		Username: "your-username",
		Password: "your-password",
	}

	// Your university domainURL
	domainURL := "https://eclass.your-university.gr"

	// Login function will return a new http.client
	client, err := auth.Login(context.Background(), domainURL, creds, nil)
	if err != nil {
		return
	}

	// Do something as authenticated user
	_, err = client.Get(domainURL + "/main/my_courses.php")
	if err != nil {
		return
	}
}
