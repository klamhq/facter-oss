package models

// User struct is defining a system user, attributes should be corresponding to /etc/passwd attribute
type User struct {
	Uid           string    `json:"uid,omitempty"`
	Gid           string    `json:"gid,omitempty"`
	Username      string    `json:"username,omitempty"`
	Name          string    `json:"name,omitempty"`
	HomeDir       string    `json:"home_dir,omitempty"`
	Session       []Session `json:"session,omitempty"`
	CanBecomeRoot bool      `json:"canbecomeroot,omitempty"`
	Shell         string    `json:"shell,omitempty"`
}

// Users struct is used to store a slice of users and offer help method like `ToProtoBuf`
type Users struct {
	Users []User `json:"users,omitempty"`
}

type Session struct {
	Connected bool   `json:"connected,omitempty"`
	Terminal  string `json:"terminal,omitempty"`
	Started   int64  `json:"started,omitempty"`
	Host      string `json:"host,omitempty"`
}
