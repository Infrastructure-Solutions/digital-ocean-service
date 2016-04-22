package usecases

import (
	"fmt"

	"github.com/Tinker-Ware/digital-ocean-service/domain"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

func (interactor DOInteractor) ShowKeys(token string) ([]domain.Key, error) {

	client := getClient(token)

	doKeys, _, err := client.Keys.List(nil)
	if err != nil {
		return nil, err
	}

	keys := []domain.Key{}
	fmt.Printf("%#v\n", doKeys)
	for _, k := range doKeys {
		key := domain.Key{
			ID:          k.ID,
			Name:        k.Name,
			Fingerprint: k.Fingerprint,
			PublicKey:   k.PublicKey,
		}
		keys = append(keys, key)
	}

	return keys, nil
}

func getClient(token string) *godo.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := godo.NewClient(tc)

	return client
}

func (interactor DOInteractor) CreateKey(name, publicKey, token string) (*domain.Key, error) {

	client := getClient(token)

	k := &godo.KeyCreateRequest{
		Name:      name,
		PublicKey: publicKey,
	}

	doKey, _, err := client.Keys.Create(k)
	if err != nil {
		return nil, err
	}

	key := domain.Key{
		ID:          doKey.ID,
		Name:        doKey.Name,
		PublicKey:   doKey.PublicKey,
		Fingerprint: doKey.Fingerprint,
	}

	return &key, nil
}
