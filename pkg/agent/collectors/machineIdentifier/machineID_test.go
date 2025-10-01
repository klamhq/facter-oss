package machineIdentifier

import (
	"testing"

	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/klamhq/facter-oss/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const machineUUID = "testdata/machineId/product_uuid"
const machineID = "testdata/machineId/machine-id"

func TestGetMachineID(t *testing.T) {
	if utils.IsRoot() {
		m := &models.MachineID{}
		m.MachineId = "2555086021244af089fa509cd3264ee7"
		m.UUID = "4c4c4544-004d-5210-8059-b8c04f565432"
		var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
		logger := factory.New(logrus.ErrorLevel)
		res, err := GetMachineID(logger, machineID, machineUUID)
		assert.Nil(t, err)
		assert.Equal(t, m.MachineId, res.MachineId)
		assert.Equal(t, m.UUID, res.UUID)
	} else {
		// If not root, we expect the machine ID to be empty
		// and the UUID to be empty as well
		m := &models.MachineID{}
		m.MachineId = "2555086021244af089fa509cd3264ee7"
		var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
		logger := factory.New(logrus.ErrorLevel)
		res, err := GetMachineID(logger, machineID, machineUUID)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	}
}

func TestFailedGetMachineId(t *testing.T) {
	if utils.IsRoot() {
		var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
		logger := factory.New(logrus.ErrorLevel)
		res, err := GetMachineID(logger, "", "")
		assert.Error(t, err)
		assert.Nil(t, res)
	}
}

func TestFailedNoRootMachineUUID(t *testing.T) {
	instance := &models.MachineID{}
	var factory utils.LoggerFactory = &utils.DefaultLoggerFactory{}
	logger := factory.New(logrus.ErrorLevel)
	uuid, err := readProductUUID("")
	assert.Equal(t, uuid, instance.UUID)
	if err == nil {
		instance.UUID = uuid
	} else {
		logger.Debug("Unable to fetch machine UUID.")
	}
}
