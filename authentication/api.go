package authentication

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

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
