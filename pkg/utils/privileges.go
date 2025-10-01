package utils

import "os"

// IsRoot return `true` if current process is runned with root privileges.
func IsRoot() bool {
	return os.Geteuid() == 0
}
