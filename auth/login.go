package auth

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Login function authenticates user to eclass and stores session to the client provided.
// After successful login, returns the http.Client.
func Login(
	ctx context.Context,
	domainURL string,
	creds Credentials,
	client *http.Client,
) (*http.Client, error) {
	if domainURL == "" || !strings.Contains(domainURL, "https://eclass") {
		return nil, errors.Errorf("invalid domain url: %s", domainURL)
	}

	if creds.UsernameEmpty() || creds.PasswordEmpty() {
		return nil, errors.Errorf("invalid credentials")
	}

	if client == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		client = &http.Client{Jar: jar}
	}

	if client.Jar == nil {
		var err error
		client.Jar, err = cookiejar.New(nil)
		if err != nil {
			return nil, errors.Errorf("could not create new cookiejar, error: %v", err)
		}
	}

	err := headHomepage(ctx, domainURL, client)
	if err != nil {
		return nil, errors.Errorf("could not head homepage, error: %v", err)
	}

	err = postLogin(ctx, domainURL, creds, client)
	if err != nil {
		return nil, errors.Errorf("could not login, error: %v", err)
	}

	return client, nil
}

func headHomepage(ctx context.Context, domainURL string, client *http.Client) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, domainURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("site unavailable: first Get")
	}

	return nil
}

func postLogin(ctx context.Context, domainURL string, creds Credentials, client *http.Client,
) error {
	form := make(url.Values, 3)
	form.Add("uname", creds.Username)
	form.Add("pass", creds.Password)
	form.Add("submit", "")

	rdr := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, domainURL, rdr)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("could not login; status code %v", res.StatusCode)
	}

	return nil
}
