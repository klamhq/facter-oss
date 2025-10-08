package main

import (
	"os"

	"github.com/klamhq/facter-oss/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
