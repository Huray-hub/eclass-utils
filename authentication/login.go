package authentication

import (
	"context"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Credentials struct {
	Username string
	Password string
}

// Login function authenticates user to eclass and stores session to the client provided
func Login(ctx context.Context, domainURL string, creds Credentials, client *http.Client) error {
	if domainURL == "" || !strings.Contains(domainURL, "https://eclass") {
		return errors.Errorf("invalid domain url: %s", domainURL)
	}
	if client.Jar == nil {
		return errors.New("client cookie jar is nil")
	}

	err := headHomepage(ctx, domainURL, client)
	if err != nil {
		return err
	}

	err = postLogin(ctx, domainURL, creds, client)
	return err
}
