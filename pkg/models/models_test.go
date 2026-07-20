package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIpInfo(t *testing.T) {
	ipInfo := IpInfo{
		Ip:        "192.168.1.1",
		Forwarded: "10.0.0.1, 172.16.0.1",
	}

	assert.Equal(t, "192.168.1.1", ipInfo.Ip)
	assert.Equal(t, "10.0.0.1, 172.16.0.1", ipInfo.Forwarded)
}

func TestGeoIpInfo(t *testing.T) {
	geoIpInfo := GeoIpInfo{
		GeoIpInfoLocationLatitude:  37.7749,
		GeoIpInfoLocationLongitude: -122.4194,
		GeoIpInfoAccuracy:          100,
	}

	assert.Equal(t, 37.7749, geoIpInfo.GeoIpInfoLocationLatitude)
	assert.Equal(t, -122.4194, geoIpInfo.GeoIpInfoLocationLongitude)
	assert.Equal(t, int32(100), geoIpInfo.GeoIpInfoAccuracy)
}

func TestGeoIp(t *testing.T) {
	geoIp := GeoIp{
		Location: Location{
			Lat: 37.7749,
			Lng: -122.4194,
		},
		Accuracy: 100,
	}

	assert.Equal(t, 37.7749, geoIp.Location.Lat)
	assert.Equal(t, -122.4194, geoIp.Location.Lng)
	assert.Equal(t, int32(100), geoIp.Accuracy)
}

func TestApplication(t *testing.T) {
	app := Application{
		Name: "test-app",
	}

	assert.Equal(t, "test-app", app.Name)

	// Test JSON marshaling
	jsonData, err := json.Marshal(app)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), "test-app")

	// Test JSON unmarshaling
	var unmarshaledApp Application
	err = json.Unmarshal(jsonData, &unmarshaledApp)
	assert.NoError(t, err)
	assert.Equal(t, app.Name, unmarshaledApp.Name)
}

func TestDockerImage(t *testing.T) {
	image := Images{
		Containers:  2,
		Created:     1234567890,
		ID:          "sha256:abc123",
		ParentID:    "sha256:parent123",
		RepoDigests: []string{"test@sha256:digest1"},
		RepoTags:    []string{"test:latest", "test:v1.0"},
		SharedSize:  1024,
		Size:        2048,
	}

	assert.Equal(t, int64(2), image.Containers)
	assert.Equal(t, "sha256:abc123", image.ID)
	assert.Len(t, image.RepoTags, 2)
	assert.Contains(t, image.RepoTags, "test:latest")

	// Test JSON marshaling
	jsonData, err := json.Marshal(image)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
}

func TestDockerContainer(t *testing.T) {
	container := Containers{
		ID:          "container123",
		Name:        "test-container",
		Image:       "nginx:latest",
		ImageID:     "sha256:image123",
		Command:     "/bin/sh",
		Created:     1234567890,
		State:       "running",
		Status:      "Up 2 hours",
		NetworkMode: "bridge",
		Ports: []ContainerPort{
			{
				IP:          "0.0.0.0",
				PrivatePort: 80,
				PublicPort:  8080,
				Type:        "tcp",
			},
		},
		Networks: []ContainerNetwork{
			{
				MacAddress:  "02:42:ac:11:00:02",
				NetworkID:   "net123",
				EndpointID:  "endpoint123",
				Gateway:     "172.17.0.1",
				IPAddress:   "172.17.0.2",
				IPPrefixLen: 16,
			},
		},
		Mounts: []ContainerMountPoint{
			{
				Type:        "volume",
				Name:        "data",
				Source:      "/var/lib/docker/volumes/data/_data",
				Destination: "/data",
				Driver:      "local",
				Mode:        "rw",
				RW:          true,
				Propagation: "rprivate",
			},
		},
	}

	assert.Equal(t, "container123", container.ID)
	assert.Equal(t, "test-container", container.Name)
	assert.Equal(t, "running", container.State)
	assert.Len(t, container.Ports, 1)
	assert.Equal(t, int32(8080), container.Ports[0].PublicPort)
	assert.Len(t, container.Networks, 1)
	assert.Equal(t, "172.17.0.2", container.Networks[0].IPAddress)
	assert.Len(t, container.Mounts, 1)
	assert.Equal(t, "data", container.Mounts[0].Name)
}

func TestDockerNetwork(t *testing.T) {
	network := DockerNetworks{
		ID:       "net123",
		Name:     "bridge",
		Driver:   "bridge",
		Scope:    "local",
		Internal: false,
	}

	assert.Equal(t, "net123", network.ID)
	assert.Equal(t, "bridge", network.Name)
	assert.Equal(t, "local", network.Scope)
	assert.False(t, network.Internal)

	// Test JSON marshaling
	jsonData, err := json.Marshal(network)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), "bridge")
}

func TestDockerStruct(t *testing.T) {
	docker := Docker{
		Images: []Images{
			{
				ID:       "img1",
				RepoTags: []string{"test:latest"},
			},
		},
		Containers: []Containers{
			{
				ID:    "cont1",
				Name:  "test",
				State: "running",
			},
		},
		Networks: []DockerNetworks{
			{
				ID:     "net1",
				Name:   "bridge",
				Driver: "bridge",
			},
		},
	}

	assert.Len(t, docker.Images, 1)
	assert.Len(t, docker.Containers, 1)
	assert.Len(t, docker.Networks, 1)
	assert.Equal(t, "img1", docker.Images[0].ID)
	assert.Equal(t, "cont1", docker.Containers[0].ID)
	assert.Equal(t, "net1", docker.Networks[0].ID)
}

func TestSchemaNode(t *testing.T) {
	node := SchemaNode{
		ID:     "node1",
		Labels: []string{"Server", "Linux"},
		Props: []PropertyInfo{
			{
				Name:      "hostname",
				Type:      "string",
				Mandatory: true,
			},
		},
	}

	assert.Equal(t, "node1", node.ID)
	assert.Len(t, node.Labels, 2)
	assert.Contains(t, node.Labels, "Server")
	assert.Len(t, node.Props, 1)
	assert.Equal(t, "hostname", node.Props[0].Name)
	assert.True(t, node.Props[0].Mandatory)
}

func TestSchemaRelationship(t *testing.T) {
	rel := SchemaRelationship{
		Type:      "RUNS_ON",
		StartNode: "app1",
		EndNode:   "server1",
		Props: []PropertyInfo{
			{
				Name:      "since",
				Type:      "timestamp",
				Mandatory: false,
			},
		},
	}

	assert.Equal(t, "RUNS_ON", rel.Type)
	assert.Equal(t, "app1", rel.StartNode)
	assert.Equal(t, "server1", rel.EndNode)
	assert.Len(t, rel.Props, 1)
	assert.False(t, rel.Props[0].Mandatory)
}

func TestSchemaVisualization(t *testing.T) {
	viz := SchemaVisualization{
		Nodes: []SchemaNode{
			{
				ID:     "node1",
				Labels: []string{"Server"},
			},
		},
		Relationships: []SchemaRelationship{
			{
				Type:      "CONNECTS_TO",
				StartNode: "node1",
				EndNode:   "node2",
			},
		},
	}

	assert.Len(t, viz.Nodes, 1)
	assert.Len(t, viz.Relationships, 1)
	assert.Equal(t, "node1", viz.Nodes[0].ID)
	assert.Equal(t, "CONNECTS_TO", viz.Relationships[0].Type)

	// Test JSON marshaling
	jsonData, err := json.Marshal(viz)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test JSON unmarshaling
	var unmarshaledViz SchemaVisualization
	err = json.Unmarshal(jsonData, &unmarshaledViz)
	assert.NoError(t, err)
	assert.Len(t, unmarshaledViz.Nodes, 1)
	assert.Equal(t, viz.Nodes[0].ID, unmarshaledViz.Nodes[0].ID)
}
