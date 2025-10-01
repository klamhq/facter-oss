package models

type Docker struct {
	Images     []Images
	Containers []Containers
	Networks   []DockerNetworks
}

type DockerNetworks struct {
	ID       string `json:"id"`       // Unique identifier of the network
	Name     string `json:"name"`     // Name of the network
	Driver   string `json:"driver"`   // Driver used for the network
	Scope    string `json:"scope"`    // Scope of the network (e.g., local, global)
	Internal bool   `json:"internal"` // Indicates if the network is internal
}

type Images struct {
	Containers  int64    `json:"containers"`  // Number of containers using this image
	Created     int64    `json:"created"`     // Creation time of the image
	ID          string   `json:"id"`          // Unique identifier of the image
	ParentID    string   `json:"parentId"`    // ID of the parent image
	RepoDigests []string `json:"repoDigests"` // List of repository digests
	RepoTags    []string `json:"repoTags"`    // List of repository tags
	SharedSize  int64    `json:"sharedSize"`  // Size shared with other images
	Size        int64    `json:"size"`        // Size of the image
}

type Containers struct {
	ID          string                `json:"id"`          // Unique identifier of the container
	Name        string                `json:"name"`        // Name of the container
	Image       string                `json:"image"`       // Name of the image used by the container
	ImageID     string                `json:"imageID"`     // ID of the image used by the container
	Command     string                `json:"command"`     // Command used to start the container
	Created     int64                 `json:"created"`     // Creation time of the container
	Ports       []ContainerPort       `json:"ports"`       // List of ports exposed by the container
	SizeRw      int64                 `json:"sizeRw"`      // Size of the read-write layer of the container
	SizeRootFs  int64                 `json:"sizeRootFs"`  // Size of the root filesystem of the container
	State       string                `json:"state"`       // Current state of the container (e.g., running, exited)
	Status      string                `json:"status"`      // Status of the container (e.g., Up, Exited)
	NetworkMode string                `json:"networkMode"` // Network mode of the container
	Networks    []ContainerNetwork    `json:"networks"`    // List of networks the container is connected to
	Mounts      []ContainerMountPoint `json:"mounts"`      // List of mount points for the container
}
type ContainerPort struct {
	IP          string `json:"ip"`
	PrivatePort int32  `json:"privatePort"`
	PublicPort  int32  `json:"publicPort"`
	Type        string `json:"type"`
}

type ContainerNetwork struct {
	MacAddress  string `json:"macAddress"`
	NetworkID   string `json:"networkID"`
	EndpointID  string `json:"endpointID"`
	Gateway     string `json:"gateway"`
	IPAddress   string `json:"ipAddress"`
	IPPrefixLen int64  `json:"ipPrefixLen"`
}

type ContainerMountPoint struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Driver      string `json:"driver"`
	Mode        string `json:"mode"`
	RW          bool   `json:"rw"`
	Propagation string `json:"propagation"`
}
