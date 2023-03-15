package auth

import (
	"context"
	"net/http"
	"time"
)

// Session token's duration (5 minutes less actually)
const TokenDuration = 25 * time.Minute

// Session function authenticates user to eclass, stores session to the client provided
// and ensures that session will not expire. Returns the http.Client to be used, in case
// no client was provided.
func Session(
	ctx context.Context,
	domainURL string,
	creds Credentials,
	client *http.Client,
) (*http.Client, error) {
	client, err := Login(ctx, domainURL, creds, client)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTimer(TokenDuration)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := refreshToken(ctx, domainURL, creds, client)
				if err != nil {
					return
				}
			}
		}
	}()

	return client, nil
}

// refreshToken function attempts to refresh an existing session. Eclass session token expires
// after 30 minutes. Every 25 minutes, refreshToken is invoked by Session function, thus 5 minutes
// remain to expire. If we successfully HEAD homepage within this time, the token gets refreshed.
// If not, we login to create a new one.
func refreshToken(
	ctx context.Context,
	domainURL string,
	creds Credentials,
	client *http.Client,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(5 * time.Minute):
			_, err := Login(ctx, domainURL, creds, client)
			if err != nil {
				return err
			}
			return nil
		default:
			err := headHomepage(ctx, domainURL, client)
			if err == nil {
				return nil
			}
			time.Sleep(10 * time.Second)
		}
	}
}
