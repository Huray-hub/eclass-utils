package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var (
	ErrNoCredentials      = errors.New("credentials not provided")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type ErrInvalidDomain struct {
	DomainURL string
}

func (e *ErrInvalidDomain) Error() string {
	return fmt.Sprintf("invalid domain URL: %s", e.DomainURL)
}

// Login function authenticates user to eclass and stores session to the client provided.
// After successful login, returns the http.Client.
func Login(
	ctx context.Context,
	domainURL string,
	creds Credentials,
	client *http.Client,
) (*http.Client, error) {
	if domainURL == "" || !strings.Contains(domainURL, "https://eclass") {
		return nil, &ErrInvalidDomain{
			DomainURL: domainURL,
		}
	}

	if creds.UsernameEmpty() || creds.PasswordEmpty() {
		return nil, ErrNoCredentials
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
			return nil, fmt.Errorf("could not create new cookiejar: %w", err)
		}
	}

	err := headHomepage(ctx, domainURL, client)
	if err != nil {
		return nil, fmt.Errorf("could not head homepage: %w", err)
	}

	err = postLogin(ctx, domainURL, creds, client)
	if err != nil {
		return nil, fmt.Errorf("could not login: %w", err)
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

	parsedURL, err := url.Parse(domainURL)
	if err != nil {
		return err
	}

	sidBefore, err := sessionID(parsedURL, client)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("could not login: status code %d", res.StatusCode)
	}

	sidAfter, err := sessionID(parsedURL, client)
	if err != nil {
		return err
	}

	if sidBefore == sidAfter {
		return ErrInvalidCredentials
	}

	return nil
}

func sessionID(url *url.URL, client *http.Client) (string, error) {
	for _, cookie := range client.Jar.Cookies(url) {
		if cookie.Name == "PHPSESSID" {
			return cookie.Value, nil
		}
	}
	return "", errors.New("session not found")
}
