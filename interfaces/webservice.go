package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Tinker-Ware/digital-ocean-service/domain"
	"github.com/Tinker-Ware/digital-ocean-service/usecases"
)

// DOInteractor defines the functions that the digital ocean interactor should perform
type DOInteractor interface {
	GetOauthURL(id, redirectURI string, scope []string) string
	GetToken(code, id, secret, redirectURL string) (*domain.DOToken, error)
	ShowKeys(token string) ([]domain.Key, error)
	CreateKey(name, publicKey, token string) (*domain.Key, error)
	CreateDroplet(droplet domain.DropletRequest, token string) (*usecases.Instance, error)
	ListDroplets(token string) ([]domain.Droplet, error)
}

// WebServiceHandler has the necessary fields for a web interface to perform its operations
type WebServiceHandler struct {
	Interactor  DOInteractor
	ID          string
	Secret      string
	RedirectURI string
	Scopes      []string
}

type instanceResponse struct {
	Instance *usecases.Instance `json:"instance"`
}

// Login is a helper method to create an OAUTH token
func (handler WebServiceHandler) Login(res http.ResponseWriter, req *http.Request) {

	url := handler.Interactor.GetOauthURL(handler.ID, handler.RedirectURI, handler.Scopes)

	htmlIndex := `<html><body>
                Log in with <a href="` + url + `">Digital Ocean</a>
                </body></html>`

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(htmlIndex))

}

// DOCallback receives the OAUTH callback from Digital Ocean
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

// ShowKeys returns all the keys of an user in the different providers
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

// CreateKey saves a key in a provider
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

// CreateDroplet creates a droplet in Digital Ocean
func (handler WebServiceHandler) CreateDroplet(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	token := req.Header.Get("token")

	decoder := json.NewDecoder(req.Body)
	dropletRequest := domain.DropletRequest{}

	err := decoder.Decode(&dropletRequest)
	if err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	instance, err := handler.Interactor.CreateDroplet(dropletRequest, token)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	resInstance := instanceResponse{
		Instance: instance,
	}
	b, _ := json.Marshal(&resInstance)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(b)
}

// ListDroplets lists all the droplets in Digital Ocean
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
