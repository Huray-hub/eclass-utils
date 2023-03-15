package auth_test

import (
	"context"
	"time"

	"github.com/Huray-hub/eclass-utils/auth"
)

func ExampleSession() {
	// Your credentials
	creds := auth.Credentials{
		Username: "your-username",
		Password: "your-password",
	}

	// Your university domainURL
	domainURL := "https://eclass.your-university.gr"

	ctx, cancelSession := context.WithCancel(context.Background())
	client, err := auth.Session(ctx, domainURL, creds, nil)
	if err != nil {
		cancelSession()
		return
	}

	// Do something as authenticated user, session will ensure that you are staying logged in
	time.Sleep(10 * time.Hour)

	_, _ = client.Get(domainURL + "/main/my_courses.php")

	// Stop session eg. before exit
	cancelSession()
}
