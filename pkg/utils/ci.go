package utils

import "os"

const CI_KEY = "CI_JOB_ID"

// IsCI return true if runned by Gitlab-ci.
func IsCI() bool {
	return len(os.Getenv(CI_KEY)) > 0
}
