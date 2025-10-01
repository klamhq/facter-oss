package machineIdentifier

import (
	"fmt"
	"os"
	"strings"

	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/klamhq/facter-oss/pkg/utils"
	"github.com/sirupsen/logrus"
)

// GetMachineID retrieves the machine ID and UUID from the specified paths.
func GetMachineID(logger *logrus.Logger, machineID string, machineUUID string) (*models.MachineID, error) {
	instance := &models.MachineID{}
	machineIdContent, err := readMachineID(machineID)
	if err == nil {
		instance.MachineId = machineIdContent
	} else {
		logger.Errorf("Unable to fetch machine ID: %s", err)
		return nil, err
	}

	uuid, err := readProductUUID(machineUUID)
	if err == nil {
		instance.UUID = uuid
	} else {
		logger.Errorf("Unable to fetch machine UUID: %s", err)
		return nil, err
	}

	return instance, nil
}

// readProductUUID reads the product UUID from the specified path.
func readProductUUID(machineUUID string) (string, error) {
	return readFileOrError(machineUUID, true)
}

// readMachineID reads the machine ID from the specified path.
func readMachineID(machineID string) (string, error) {
	return readFileOrError(machineID, true)
}

// readFileOrError reads the content of a file at the specified path.
func readFileOrError(path string, MustAdmin bool) (string, error) {

	if MustAdmin && !utils.IsRoot() {
		return "", fmt.Errorf("must be run as root to read file %s", path)
	}
	content, err := os.ReadFile(path)
	return strings.ReplaceAll(string(content), "\n", ""), err
}
