package usecases

import (
	"github.com/Tinker-Ware/digital-ocean-service/domain"
	"github.com/digitalocean/godo"
	"github.com/jinzhu/copier"
)

type Instance struct {
	Provider string `json:"provider"`
	domain.Droplet
}

func (interactor DOInteractor) CreateDroplet(droplet domain.DropletRequest, token string) (*Instance, error) {
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

	inst := &Instance{
		Provider: "digital_ocean",
		Droplet: domain.Droplet{
			ID:                drop.ID,
			Name:              droplet.Name,
			Region:            droplet.Region,
			OperatingSystem:   drop.Image.Slug,
			PrivateNetworking: false,
			InstanceName:      drop.Size.Slug,
		},
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
	networks := domain.Networks{
		V4: networksV4,
		V6: networksV6,
	}

	inst.Networks = networks
	inst.SSHKeys = droplet.SSHKeys
	return inst, nil

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
			Name:         drops.Name,
			Region:       drops.Region.String(),
			InstanceName: drops.Size.String(),
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
