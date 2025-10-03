package application

import (
	"context"
	"os"
	"testing"

	"github.com/docker/docker/client"
	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetDockerInfo(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	docker, err := GetDockerInfo(ctx, logger)
	assert.NoError(t, err)
	assert.NotEmpty(t, docker)
}

func TestListDockerImgs_NoCli(t *testing.T) {
	os.Setenv("DOCKER_HOST", "http://fake")
	ctx := context.Background()
	logger := logrus.New()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	assert.NoError(t, err)
	err = listDockerImg(ctx, cli, &models.Docker{}, logger)
	assert.Error(t, err)
}

func TestListDockerContainers_NoCli(t *testing.T) {
	os.Setenv("DOCKER_HOST", "http://fake")
	ctx := context.Background()
	logger := logrus.New()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	assert.NoError(t, err)
	err = listContainers(ctx, cli, &models.Docker{}, logger)
	assert.Error(t, err)
}

func TestListDockerNetworks_NoCli(t *testing.T) {
	os.Setenv("DOCKER_HOST", "http://fake")
	ctx := context.Background()
	logger := logrus.New()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	assert.NoError(t, err)
	err = listNetworks(ctx, cli, &models.Docker{}, logger)
	assert.Error(t, err)
}
