package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/digital-ocean-service/domain"
)

type DOInteractor interface {
	GetOauthURL(id, redirectURI, scope string) string
	GetToken(code, id, secret, redirectURL string) (*domain.DOToken, error)
}

type WebServiceHandler struct {
	Interactor  DOInteractor
	ID          string
	Secret      string
	RedirectURI string
}

func (handler WebServiceHandler) Login(res http.ResponseWriter, req *http.Request) {

	url := handler.Interactor.GetOauthURL(handler.ID, handler.RedirectURI, "read write")

	htmlIndex := `<html><body>
                Log in with <a href="` + url + `">Digital Ocean</a>
                </body></html>`

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(htmlIndex))

}

func (handler WebServiceHandler) DOCallback(res http.ResponseWriter, req *http.Request) {
	code := req.FormValue("code")

	token, err := handler.Interactor.GetToken(code, handler.ID, handler.Secret, handler.RedirectURI)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(token)
	if err != nil {
		fmt.Println(err.Error())

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write(b)

}
