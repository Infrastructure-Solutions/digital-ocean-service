package usecases

import "net/url"

type DOInteractor struct {
}

const authURL = "https://cloud.digitalocean.com/v1/oauth/authorize"

func (interactor DOInteractor) GetOauthURL(id, redirectURI, scope string) string {
	u, _ := url.Parse(authURL)
	q := u.Query()
	q.Set("client_id", id)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", scope)
	q.Set("response_type", "code")

	u.RawQuery = q.Encode()

	return u.String()
}
