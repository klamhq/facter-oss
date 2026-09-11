package utils

import (
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
)

// GetHostnameFromInventory extracts the hostname from the given inventory message.
func GetHostnameFromInventory(inv *schema.InventoryRequest) string {
	if inv == nil {
		return ""
	}
	switch v := inv.Content.(type) {
	case *schema.InventoryRequest_Full:
		if v.Full != nil {
			return v.Full.Hostname
		}
	case *schema.InventoryRequest_Delta:
		if v.Delta != nil {
			return v.Delta.Hostname
		}
	case *schema.InventoryRequest_Revision:
		if v.Revision != nil {
			return v.Revision.Hostname
		}
	}
	return ""
}
