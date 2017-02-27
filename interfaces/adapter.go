package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	jose "github.com/dvsekhvalnov/jose2go"
)

// Adapter is the signature of an HTTPHandler for middlewares
type Adapter func(http.Handler) http.Handler

// Adapt takes several Adapters and calls them in order
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

type integrationsWrapper struct {
	Integrations []integration `json:"integrations"`
}

const integrationsURL string = "/api/v1/users/%s/integration"

// GetToken gets the token from the users microservice
func GetToken(apiURL string, salt string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userToken := r.Header.Get("authorization")

			sp := strings.Split(userToken, " ")
			payload, _, _ := jose.Decode(sp[1], []byte(salt))

			var objmap map[string]*json.RawMessage
			json.Unmarshal([]byte(payload), &objmap)

			var aud string
			json.Unmarshal(*objmap["aud"], &aud)

			vals := strings.Split(aud, ":")

			path := fmt.Sprintf(integrationsURL, vals[1])

			request, _ := http.NewRequest(http.MethodGet, apiURL+path, nil)
			request.Header.Set("authorization", userToken)

			client := &http.Client{}
			resp, _ := client.Do(request)
			defer resp.Body.Close()

			integrations := integrationsWrapper{}
			decoder := json.NewDecoder(resp.Body)
			decoder.Decode(&integrations)
			var token string

			for _, integ := range integrations.Integrations {
				if integ.Provider == "digital_ocean" {
					token = integ.Token
				}
			}

			r.Header.Set("key", token)

			h.ServeHTTP(w, r)
		})
	}
}
