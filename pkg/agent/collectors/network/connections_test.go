package network

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConnections(t *testing.T) {
	conn, err := Connections(logrus.New())
	assert.NoError(t, err)
	assert.NotEmpty(t, conn)
}

func TestGetConnections(t *testing.T) {
	conn, err := getConnections(logrus.New())
	assert.NoError(t, err)
	assert.NotEmpty(t, conn)
}
