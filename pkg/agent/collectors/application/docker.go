package application

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/sirupsen/logrus"
)

func GetDockerInfo(ctx context.Context, logger *logrus.Logger) (*models.Docker, error) {
	docker := &models.Docker{}

	cli, err := getDockerClient()
	if err != nil {
		return nil, err
	}

	if err := listDockerImg(ctx, cli, docker, logger); err != nil {
		return nil, err
	}

	if err := listContainers(ctx, cli, docker, logger); err != nil {
		return nil, err
	}

	if err := listNetworks(ctx, cli, docker, logger); err != nil {
		return nil, err
	}

	return docker, nil
}

// getDockerClient initializes an official Docker client from env
func getDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// listDockerImg retrieves Docker images using the official SDK
func listDockerImg(ctx context.Context, cli *client.Client, docker *models.Docker, logger *logrus.Logger) error {
	imgs, err := cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		logger.Errorf("Error during listing of docker images %s", err)
		return err
	}
	docker.Images = []models.Images{}
	for _, img := range imgs {
		imgModel := models.Images{
			Containers:  img.Containers,
			Created:     img.Created,
			ID:          img.ID,
			ParentID:    img.ParentID,
			RepoDigests: img.RepoDigests,
			RepoTags:    img.RepoTags,
			SharedSize:  img.SharedSize,
			Size:        img.Size,
		}
		docker.Images = append(docker.Images, imgModel)
	}
	return nil
}

// listContainers retrieves Docker containers using the official SDK
func listContainers(ctx context.Context, cli *client.Client, docker *models.Docker, logger *logrus.Logger) error {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		logger.Errorf("Error during listing of docker containers %s", err)
		return err
	}
	docker.Containers = []models.Containers{}
	for _, container := range containers {
		containerPort := []models.ContainerPort{}
		for _, ports := range container.Ports {
			containerPort = append(containerPort, models.ContainerPort{
				IP:          ports.IP,
				PrivatePort: int32(ports.PrivatePort),
				PublicPort:  int32(ports.PublicPort),
				Type:        ports.Type,
			})
		}
		containerNetwork := []models.ContainerNetwork{}
		for _, networks := range container.NetworkSettings.Networks {
			containerNetwork = append(containerNetwork, models.ContainerNetwork{
				MacAddress:  networks.MacAddress,
				NetworkID:   networks.NetworkID,
				EndpointID:  networks.EndpointID,
				IPAddress:   networks.IPAddress,
				IPPrefixLen: int64(networks.IPPrefixLen),
				Gateway:     networks.Gateway,
			})
		}
		containerMountPoint := []models.ContainerMountPoint{}
		for _, mount := range container.Mounts {
			containerMountPoint = append(containerMountPoint, models.ContainerMountPoint{
				Type:        string(mount.Type),
				Name:        mount.Name,
				Source:      mount.Source,
				Destination: mount.Destination,
				Mode:        mount.Mode,
				RW:          mount.RW,
				Driver:      mount.Driver,
				Propagation: string(mount.Propagation),
			})
		}
		cModel := models.Containers{
			ID:          container.ID,
			Name:        string(container.Names[0]),
			Image:       container.Image,
			ImageID:     container.ImageID,
			Command:     container.Command,
			Created:     container.Created,
			Ports:       containerPort,
			SizeRw:      container.SizeRw,
			SizeRootFs:  container.SizeRootFs,
			Status:      container.Status,
			State:       container.State,
			NetworkMode: container.HostConfig.NetworkMode,
			Networks:    containerNetwork,
			Mounts:      containerMountPoint,
		}

		docker.Containers = append(docker.Containers, cModel)
	}
	return nil
}

// listNetworks retrieves Docker networks using the official SDK
func listNetworks(ctx context.Context, cli *client.Client, docker *models.Docker, logger *logrus.Logger) error {
	networks, err := cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		logger.Errorf("Error during listing of docker networks %s", err)
		return err
	}
	docker.Networks = []models.DockerNetworks{}
	for _, network := range networks {
		nModel := models.DockerNetworks{
			ID:       network.ID,
			Name:     network.Name,
			Driver:   network.Driver,
			Scope:    network.Scope,
			Internal: network.Internal,
		}
		docker.Networks = append(docker.Networks, nModel)
	}
	return nil
}
