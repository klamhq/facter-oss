package utils

import (
	"testing"

	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetHostnameFromInventory_Full(t *testing.T) {
	expected := "host-full"
	inv := &schema.InventoryRequest{
		Content: &schema.InventoryRequest_Full{
			Full: &schema.HostInventory{
				Hostname: expected,
			},
		},
	}
	got := GetHostnameFromInventory(inv)
	assert.Equal(t, expected, got, "Hostnames should match")
}

func TestGetHostnameFromInventory_Delta(t *testing.T) {
	expected := "host-delta"
	inv := &schema.InventoryRequest{
		Content: &schema.InventoryRequest_Delta{
			Delta: &schema.HostDeltaInventory{
				Hostname: expected,
			},
		},
	}
	got := GetHostnameFromInventory(inv)
	assert.Equal(t, expected, got, "Hostnames should match")
}

func TestGetHostnameFromInventory_UnknownType(t *testing.T) {
	inv := &schema.InventoryRequest{
		Content: nil,
	}
	got := GetHostnameFromInventory(inv)
	assert.Empty(t, got, "Hostname should be empty for unknown inventory type")
}
