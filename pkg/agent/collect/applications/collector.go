package applications

import (
	"context"
	"fmt"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/application"
	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
)

type ApplicationsCollectorImpl struct {
	log *logrus.Logger
	cfg *options.ApplicationsOptions
}

func New(log *logrus.Logger, cfg *options.ApplicationsOptions) *ApplicationsCollectorImpl {

	return &ApplicationsCollectorImpl{
		log: log,
		cfg: cfg,
	}
}

func (c *ApplicationsCollectorImpl) Name() string { return "applications" }

func (c *ApplicationsCollectorImpl) CollectApplications(ctx context.Context) ([]*schema.Application, error) {
	c.log.Info("Crafting applications")

	if !c.cfg.Docker.Enabled {
		return nil, nil
	}

	info, err := application.GetDockerInfo(ctx, c.log)
	if err != nil {
		return nil, fmt.Errorf("get docker info: %w", err)
	}

	docker := &schema.Docker{
		Images:     make([]*schema.ContainersImages, 0, len(info.Images)),
		Networks:   make([]*schema.DockerNetworks, 0, len(info.Networks)),
		Containers: make([]*schema.Containers, 0, len(info.Containers)),
	}

	// Images
	for _, img := range info.Images {
		docker.Images = append(docker.Images, &schema.ContainersImages{
			Id:         img.ID,
			RepoTags:   img.RepoTags,
			Created:    img.Created,
			Size:       img.Size,
			SharedSize: img.SharedSize,
			ParentId:   img.ParentID,
		})
	}

	// Networks
	for _, net := range info.Networks {
		docker.Networks = append(docker.Networks, &schema.DockerNetworks{
			Id:       net.ID,
			Name:     net.Name,
			Driver:   net.Driver,
			Internal: net.Internal,
			Scope:    net.Scope,
		})
	}

	// Containers
	for _, ctr := range info.Containers {
		// mounts
		mounts := make([]*schema.ContainerMounts, 0, len(ctr.Mounts))
		for _, m := range ctr.Mounts {
			mounts = append(mounts, &schema.ContainerMounts{
				Name:        m.Name,
				Source:      m.Source,
				Destination: m.Destination,
				Driver:      m.Driver,
				Mode:        m.Mode,
				Rw:          m.RW,
				Propagation: m.Propagation,
				Type:        m.Type,
			})
		}

		ports := make([]*schema.ContainerPorts, 0, len(ctr.Ports))
		for _, p := range ctr.Ports {
			ports = append(ports, &schema.ContainerPorts{
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Type:        p.Type,
				Ip:          p.IP,
			})
		}

		nets := make([]*schema.ContainerNetworks, 0, len(ctr.Networks))
		for _, n := range ctr.Networks {
			nets = append(nets, &schema.ContainerNetworks{
				MacAddress:  n.MacAddress,
				IpAddress:   n.IPAddress,
				Gateway:     n.Gateway,
				EndpointId:  n.EndpointID,
				NetworkId:   n.NetworkID,
				IpPrefixLen: n.IPPrefixLen,
			})
		}

		docker.Containers = append(docker.Containers, &schema.Containers{
			Id:          ctr.ID,
			Name:        ctr.Name,
			Created:     ctr.Created,
			Image:       ctr.Image,
			ImageId:     ctr.ImageID,
			Command:     ctr.Command,
			SizeRootFs:  ctr.SizeRootFs,
			Mounts:      mounts,
			State:       ctr.State,
			Status:      ctr.Status,
			Ports:       ports,
			Networks:    nets,
			NetworkMode: ctr.NetworkMode,
			SizeRw:      ctr.SizeRw,
		})
	}

	// on ne pousse une Application que si Docker est actif
	app := &schema.Application{Docker: docker}
	return []*schema.Application{app}, nil
}
