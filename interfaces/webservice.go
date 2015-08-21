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
	ShowKeys(token string) ([]domain.Key, error)
	CreateKey(name, publicKey, token string) (*domain.Key, error)
	CreateDroplet(droplet domain.DropletRequest, token string) (*domain.Droplet, error)
	ListDroplets(token string) ([]domain.Droplet, error)
}

type UserRepo interface {
}

type WebServiceHandler struct {
	Interactor  DOInteractor
	UserRepo    UserRepo
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

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(b)

}

func (handler WebServiceHandler) ShowKeys(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")
	keys, err := handler.Interactor.ShowKeys(token)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	keysB, _ := json.Marshal(keys)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(keysB)

}

func (handler WebServiceHandler) CreateKey(res http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	token := req.Header.Get("token")

	key := &domain.Key{}

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(key)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	key, err = handler.Interactor.CreateKey(key.Name, key.PublicKey, token)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	keyB, _ := json.Marshal(key)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(keyB)
}

func (handler WebServiceHandler) CreateDroplet(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	token := req.Header.Get("token")

	decoder := json.NewDecoder(req.Body)
	dropletRequest := domain.DropletRequest{}

	err := decoder.Decode(&dropletRequest)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	droplet, err := handler.Interactor.CreateDroplet(dropletRequest, token)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(droplet)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(b)
}

func (handler WebServiceHandler) ListDroplets(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("token")

	droplets, err := handler.Interactor.ListDroplets(token)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	dB, _ := json.Marshal(droplets)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(dB)
}
