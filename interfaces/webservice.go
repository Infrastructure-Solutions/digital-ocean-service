package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Tinker-Ware/digital-ocean-service/domain"
	"github.com/Tinker-Ware/digital-ocean-service/usecases"
	"github.com/gorilla/mux"
)

type DOInteractor interface {
	GetOauthURL(id, redirectURI string, scope []string) string
	GetToken(code, id, secret, redirectURL string) (*domain.DOToken, error)
	ShowKeys(token string) ([]domain.Key, error)
	CreateKey(name, publicKey, token string) (*domain.Key, error)
	CreateDroplet(droplet domain.DropletRequest, token string) (*usecases.Instance, error)
	ListDroplets(token string) ([]domain.Droplet, error)
	GetDroplet(id int, token string) (*usecases.Instance, error)
}

type UserRepo interface {
}

type WebServiceHandler struct {
	Interactor  DOInteractor
	UserRepo    UserRepo
	ID          string
	Secret      string
	RedirectURI string
	Scopes      []string
}

type instanceResponse struct {
	Instance *usecases.Instance `json:"instance"`
}

func (handler WebServiceHandler) Login(res http.ResponseWriter, req *http.Request) {

	url := handler.Interactor.GetOauthURL(handler.ID, handler.RedirectURI, handler.Scopes)

	htmlIndex := `<html><body>
                Log in with <a href="` + url + `">Digital Ocean</a>
                </body></html>`

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(htmlIndex))

}

const providerToken string = "provider-token"

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
	token := req.Header.Get(providerToken)
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

	token := req.Header.Get(providerToken)

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

	token := req.Header.Get(providerToken)

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

func (handler WebServiceHandler) ListDroplets(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get(providerToken)

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

// GetInstance gets an instance from any provider
func (handler WebServiceHandler) GetInstance(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get(providerToken)
	vars := mux.Vars(req)
	id := vars["instanceID"]
	instanceID, _ := strconv.Atoi(id)

	droplet, err := handler.Interactor.GetDroplet(instanceID, token)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	response := instanceResponse{
		Instance: droplet,
	}

	dB, _ := json.Marshal(response)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(dB)

}
