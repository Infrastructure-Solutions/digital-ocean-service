package interfaces

import "net/http"

type DOInteractor interface {
	GetOauthURL(id, redirectURI, scope string) string
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
