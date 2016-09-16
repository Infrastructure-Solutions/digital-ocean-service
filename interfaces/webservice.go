package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Tinker-Ware/digital-ocean-service/domain"
	"github.com/Tinker-Ware/digital-ocean-service/usecases"
	"github.com/gorilla/mux"
)

// DOInteractor defines the functions that the digital ocean interactor should perform
type DOInteractor interface {
	GetOauthURL(id, redirectURI string, scope []string) string
	GetToken(code, id, secret, redirectURL string) (*domain.DOToken, error)
	ShowKeys(token string) ([]domain.Key, error)
	CreateKey(name, publicKey, token string) (*domain.Key, error)
	CreateDroplet(droplet domain.DropletRequest, token string) (*usecases.Instance, error)
	ListDroplets(token string) ([]domain.Droplet, error)
	GetDroplet(id int, token string) (*usecases.Instance, error)
}

// WebServiceHandler has the necessary fields for a web interface to perform its operations
type WebServiceHandler struct {
	Interactor  DOInteractor
	ID          string
	Secret      string
	RedirectURI string
	Scopes      []string
	APIHost     string
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

const providerToken string = "provider-token"

type oauthWrapper struct {
	OauthRequest oauthRequest `json:"oauth_request"`
}

type oauthRequest struct {
	UserID int    `json:"user_id"`
	Code   string `json:"code"`
	State  string `json:"state"`
}

type integrationWrapper struct {
	Integration integration `json:"integration"`
}

type integration struct {
	UserID     int    `json:"user_id"`
	Token      string `json:"token"`
	Username   string `json:"username"`
	Provider   string `json:"provider"`
	ExpireDate int    `json:"expire_date"`
}

type callbackResponse struct {
	Callback callback `json:"callback"`
}

type callback struct {
	Provider string `json:"provider"`
	Username string `json:"username"`
}

const integrationURL string = "/api/v1/users/%d/integration"

// DOCallback receives the OAUTH callback from Digital Ocean
func (handler WebServiceHandler) DOCallback(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	userToken := req.Header.Get("Authorization")

	decoder := json.NewDecoder(req.Body)

	var oauthwrapper oauthWrapper

	err := decoder.Decode(&oauthwrapper)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := handler.Interactor.GetToken(oauthwrapper.OauthRequest.Code, handler.ID, handler.Secret, handler.RedirectURI)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	wrapper := integrationWrapper{
		Integration: integration{
			UserID:     oauthwrapper.OauthRequest.UserID,
			Token:      token.AccessToken,
			Provider:   "digital_ocean",
			Username:   token.Info.Name,
			ExpireDate: token.ExpiresIn,
		},
	}

	reqBytes, _ := json.Marshal(&wrapper)

	buf := bytes.NewBuffer(reqBytes)

	path := fmt.Sprintf(integrationURL, oauthwrapper.OauthRequest.UserID)

	request, _ := http.NewRequest(http.MethodPost, handler.APIHost+path, buf)
	request.Header.Add("Authorization", userToken)
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, _ := client.Do(request)
	if resp.StatusCode != http.StatusCreated {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := callbackResponse{
		Callback: callback{
			Provider: "github",
			Username: token.Info.Name,
		},
	}

	respBytes, _ := json.Marshal(&response)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(respBytes)

}

// ShowKeys returns all the keys of an user in the different providers
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

// CreateKey saves a key in a provider
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

// CreateDroplet creates a droplet in Digital Ocean
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

// ListDroplets lists all the droplets in Digital Ocean
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
