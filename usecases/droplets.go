package usecases

import (
	"github.com/digital-ocean-service/domain"
	"github.com/digitalocean/godo"
	"github.com/jinzhu/copier"
)

func (interactor DOInteractor) CreateDroplet(droplet domain.DropletRequest, token string) (*domain.Droplet, error) {
	client := getClient(token)

	dropletRequest := &godo.DropletCreateRequest{
		Name:              droplet.Name,
		Region:            droplet.Region,
		Size:              droplet.Size,
		Backups:           droplet.Backups,
		IPv6:              droplet.IPv6,
		PrivateNetworking: droplet.PrivateNetworking,
		UserData:          droplet.UserData,
		Image: godo.DropletCreateImage{
			Slug: droplet.Image,
		},
	}

	sshkeys := []godo.DropletCreateSSHKey{}
	for _, key := range droplet.SSHKeys {
		k := godo.DropletCreateSSHKey{
			ID:          key.ID,
			Fingerprint: key.Fingerprint,
		}
		sshkeys = append(sshkeys, k)
	}
	dropletRequest.SSHKeys = sshkeys

	drop, _, err := client.Droplets.Create(dropletRequest)
	if err != nil {
		return nil, err
	}

	drp := &domain.Droplet{
		Name:     droplet.Name,
		Region:   droplet.Region,
		Size:     droplet.Size,
		UserData: droplet.UserData,
	}

	networksV4 := []domain.NetworkV4{}
	for _, net := range drop.Networks.V4 {
		n := domain.NetworkV4{}
		copier.Copy(n, net)
		networksV4 = append(networksV4, n)
	}

	networksV6 := []domain.NetworkV6{}
	for _, net := range drop.Networks.V6 {
		n := domain.NetworkV6{}
		copier.Copy(n, net)
		networksV6 = append(networksV6, n)
	}

	return drp, nil

}

func (interactor DOInteractor) ListDroplets(token string) ([]domain.Droplet, error) {

	client := getClient(token)

	doDrops, _, err := client.Droplets.List(nil)
	if err != nil {
		return nil, err
	}
	droplets := []domain.Droplet{}

	for _, drops := range doDrops {
		drp := domain.Droplet{
			Name:   drops.Name,
			Region: drops.Region.String(),
			Size:   drops.Size.String(),
		}

		networksV4 := []domain.NetworkV4{}
		for _, net := range drops.Networks.V4 {
			n := domain.NetworkV4{}
			copier.Copy(n, net)
			networksV4 = append(networksV4, n)
		}

		networksV6 := []domain.NetworkV6{}
		for _, net := range drops.Networks.V6 {
			n := domain.NetworkV6{}
			copier.Copy(n, net)
			networksV6 = append(networksV6, n)
		}

		droplets = append(droplets, drp)

	}

	return droplets, nil
}
