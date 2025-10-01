package utils

import (
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

// GetHostnameFromInventory extracts the hostname from the given inventory message.
func GetHostnameFromInventory(inv *schema.InventoryRequest) string {
	switch v := inv.Content.(type) {
	case *schema.InventoryRequest_Full:
		return v.Full.Hostname
	case *schema.InventoryRequest_Delta:
		return v.Delta.Hostname
	default:
		return ""
	}
}
