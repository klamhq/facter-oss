package inventory

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"time"

	"github.com/google/go-cmp/cmp"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

// StableHash returns a deterministic hash for a protobuf message.
// It is used to detect content changes on the same logical item.
func StableHash(msg proto.Message) uint64 {
	if msg == nil {
		return 0
	}

	b, _ := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(msg)

	h := fnv.New64a()
	_, _ = h.Write(b)
	return h.Sum64()
}

// DebugProtoDiff returns a readable diff for logs.
func DebugProtoDiff(oldMsg, newMsg proto.Message) string {
	if oldMsg == nil && newMsg == nil {
		return "No difference"
	}

	diff := cmp.Diff(
		oldMsg,
		newMsg,
		protocmp.Transform(),
		protocmp.IgnoreFields(oldMsg),
	)
	if diff == "" {
		return "No difference"
	}
	return fmt.Sprintf("Proto diff:\n%s", diff)
}

func isNilProtoValue[T proto.Message](v T) bool {
	if reflect.TypeOf(v) == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Func:
		return rv.IsNil()
	default:
		return false
	}
}

// sameProtoMessage compares two protobuf messages.
func sameProtoMessage[T proto.Message](oldMsg, newMsg T) bool {
	if isNilProtoValue(oldMsg) || isNilProtoValue(newMsg) {
		return isNilProtoValue(oldMsg) && isNilProtoValue(newMsg)
	}
	return proto.Equal(oldMsg, newMsg)
}

// DiffGenericByHash compares repeated protobuf entries by identity and stable hash.
// Same identity + different hash => modified item.
func DiffGenericByHash[T proto.Message](
	oldList, newList []T,
	getKey func(T) string,
) (added, removed, changed []T) {
	oldByKey := make(map[string]T, len(oldList))
	for _, item := range oldList {
		if isNilProtoValue(item) {
			continue
		}
		oldByKey[getKey(item)] = item
	}

	for _, item := range newList {
		if isNilProtoValue(item) {
			continue
		}

		key := getKey(item)
		oldItem, exists := oldByKey[key]
		if !exists {
			added = append(added, item)
			continue
		}

		if StableHash(oldItem) != StableHash(item) {
			changed = append(changed, item)
		}

		delete(oldByKey, key)
	}

	for _, item := range oldByKey {
		if isNilProtoValue(item) {
			continue
		}
		removed = append(removed, item)
	}

	return added, removed, changed
}

// Identity helpers for each entity.
func packageKey(item *schema.Package) string {
	if item == nil {
		return ""
	}
	return item.GetName() + ":" + item.GetVersion() + ":" + item.GetArchitecture()
}

func userKey(item *schema.User) string {
	if item == nil {
		return ""
	}
	return item.GetUsername()
}

func processKey(item *schema.Process) string {
	if item == nil {
		return ""
	}
	return fmt.Sprintf("%d", item.GetPid())
}

func systemdServiceKey(item *schema.SystemdService) string {
	if item == nil {
		return ""
	}
	return item.GetName()
}

func knownHostKey(item *schema.KnownHost) string {
	if item == nil {
		return ""
	}
	return item.GetHostname() + ":" + item.GetFingerprint()
}

func sshKeyInfoKey(item *schema.SshKeyInfo) string {
	if item == nil {
		return ""
	}
	return item.GetFingerprint()
}

func sshKeyAccessKey(item *schema.SshKeyAccess) string {
	if item == nil {
		return ""
	}
	return item.GetFingerprint() + ":" + item.GetAsUser()
}

func applicationKey(item *schema.Application) string {
	if item == nil {
		return ""
	}
	return fmt.Sprintf("%d", StableHash(item))
}

// ComputeDelta compares two HostInventory values and returns only the changed pieces.
func ComputeDelta(oldInv, newInv *schema.HostInventory, logger *logrus.Logger) *schema.HostDeltaInventory {
	if oldInv == nil {
		oldInv = &schema.HostInventory{}
	}
	if newInv == nil {
		newInv = &schema.HostInventory{}
	}

	delta := &schema.HostDeltaInventory{
		Hostname:  newInv.GetHostname(),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// Singular protobuf messages.
	if !sameProtoMessage(oldInv.GetPlatform(), newInv.GetPlatform()) {
		logger.Debugf("Platform changed: %s", DebugProtoDiff(oldInv.GetPlatform(), newInv.GetPlatform()))
		delta.Platform = newInv.GetPlatform()
	}

	if !sameProtoMessage(oldInv.GetNetwork(), newInv.GetNetwork()) {
		logger.Debugf("Network changed: %s", DebugProtoDiff(oldInv.GetNetwork(), newInv.GetNetwork()))
		delta.Network = newInv.GetNetwork()
	}

	// Repeated entities.
	delta.PackagesAdded, delta.PackagesRemoved, _ = DiffGenericByHash(
		oldInv.GetPackages(),
		newInv.GetPackages(),
		func(p *schema.Package) string { return packageKey(p) },
	)

	delta.UsersAdded, delta.UsersRemoved, _ = DiffGenericByHash(
		oldInv.GetUsers(),
		newInv.GetUsers(),
		func(u *schema.User) string { return userKey(u) },
	)

	delta.ProcessesAdded, delta.ProcessesRemoved, _ = DiffGenericByHash(
		oldInv.GetProcesses(),
		newInv.GetProcesses(),
		func(p *schema.Process) string { return processKey(p) },
	)

	delta.SystemdservicesAdded, delta.SystemdservicesRemoved, _ = DiffGenericByHash(
		oldInv.GetSystemdService(),
		newInv.GetSystemdService(),
		func(s *schema.SystemdService) string { return systemdServiceKey(s) },
	)

	delta.KnownhostsAdded, delta.KnownhostsRemoved, _ = DiffGenericByHash(
		oldInv.GetKnownHost(),
		newInv.GetKnownHost(),
		func(k *schema.KnownHost) string { return knownHostKey(k) },
	)

	delta.SshkeyinfoAdded, delta.SshkeyinfoRemoved, _ = DiffGenericByHash(
		oldInv.GetSshKeyInfo(),
		newInv.GetSshKeyInfo(),
		func(k *schema.SshKeyInfo) string { return sshKeyInfoKey(k) },
	)

	delta.SshkeyaccessAdded, delta.SshkeyaccessRemoved, _ = DiffGenericByHash(
		oldInv.GetSshKeyAccess(),
		newInv.GetSshKeyAccess(),
		func(k *schema.SshKeyAccess) string { return sshKeyAccessKey(k) },
	)

	delta.ApplicationsAdded, delta.ApplicationsRemoved, _ = DiffGenericByHash(
		oldInv.GetApplication(),
		newInv.GetApplication(),
		func(a *schema.Application) string { return applicationKey(a) },
	)

	return delta
}

// IsDeltaEmpty returns true when no meaningful change is present.
func IsDeltaEmpty(d *schema.HostDeltaInventory) bool {
	if d == nil {
		return true
	}

	return len(d.GetPackagesAdded()) == 0 &&
		len(d.GetPackagesRemoved()) == 0 &&
		len(d.GetUsersAdded()) == 0 &&
		len(d.GetUsersRemoved()) == 0 &&
		len(d.GetApplicationsAdded()) == 0 &&
		len(d.GetApplicationsRemoved()) == 0 &&
		len(d.GetSystemdservicesAdded()) == 0 &&
		len(d.GetSystemdservicesRemoved()) == 0 &&
		len(d.GetKnownhostsAdded()) == 0 &&
		len(d.GetKnownhostsRemoved()) == 0 &&
		len(d.GetSshkeyaccessAdded()) == 0 &&
		len(d.GetSshkeyaccessRemoved()) == 0 &&
		len(d.GetSshkeyinfoAdded()) == 0 &&
		len(d.GetSshkeyinfoRemoved()) == 0 &&
		len(d.GetProcessesAdded()) == 0 &&
		len(d.GetProcessesRemoved()) == 0 &&
		d.GetPlatform() == nil &&
		d.GetNetwork() == nil
}
