package inventory

import (
	"fmt"
	"hash/fnv"
	"time"

	"github.com/google/go-cmp/cmp"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

// StableHash return determinist hash from protobuf message
func StableHash(msg proto.Message) uint64 {
	b, _ := protojson.MarshalOptions{
		EmitUnpopulated: true, // Inclure les champs "vides"
		UseProtoNames:   true, // Noms proto stables
	}.Marshal(msg)

	h := fnv.New64a()
	h.Write(b)
	return h.Sum64()
}

// DebugProtoDiff prints the difference between two proto messages
func DebugProtoDiff(oldMsg, newMsg proto.Message) string {
	diff := cmp.Diff(
		oldMsg, newMsg,
		protocmp.Transform(), // Ignore nil/empty normalisation
		protocmp.IgnoreFields(oldMsg /* à compléter si tu veux ignorer certains champs */),
	)
	if diff == "" {
		return "No difference"
	}
	return fmt.Sprintf("Proto diff:\n%s", diff)
}

// ComputeDelta computes the difference between two HostInventory objects.
// It returns a HostDeltaInventory containing the changes.
// The delta includes added and removed packages, users, and other entities.
// It uses the hashMessage function to compare the objects efficiently.
func ComputeDelta(oldInv, newInv *schema.HostInventory, logger *logrus.Logger) *schema.HostDeltaInventory {
	delta := &schema.HostDeltaInventory{
		Hostname:  newInv.Hostname,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if oldInv.Platform == nil || !proto.Equal(oldInv.Platform, newInv.Platform) {
		logger.Debugf("Platform changed: %s", DebugProtoDiff(oldInv.Platform, newInv.Platform))
		delta.Platform = newInv.Platform
	}

	if oldInv.Network == nil || !proto.Equal(oldInv.Network, newInv.Network) {
		logger.Debugf("Network changed: %s", DebugProtoDiff(oldInv.Network, newInv.Network))
		delta.Network = newInv.Network
	}

	delta.ProcessesAdded, delta.ProcessesRemoved, _ = DiffGenericByHash(
		oldInv.Processes,
		newInv.Processes,
		func(p *schema.Process) string { return fmt.Sprintf("%d", p.Pid) },
	)

	delta.PackagesAdded, delta.PackagesRemoved, _ = DiffGenericByHash(
		oldInv.Packages,
		newInv.Packages,
		func(p *schema.Package) string { return p.Name + p.Version },
	)

	delta.UsersAdded, delta.UsersRemoved, _ = DiffGenericByHash(
		oldInv.Users,
		newInv.Users,
		func(u *schema.User) string { return u.Username },
	)
	delta.SystemdservicesAdded, delta.SystemdservicesRemoved, _ = DiffGenericByHash(
		oldInv.SystemdService,
		newInv.SystemdService,
		func(s *schema.SystemdService) string { return s.Name },
	)

	delta.KnownhostsAdded, delta.KnownhostsRemoved, _ = DiffGenericByHash(
		oldInv.KnownHost,
		newInv.KnownHost,
		func(k *schema.KnownHost) string { return k.Hostname + k.Fingerprint },
	)

	delta.SshkeyaccessAdded, delta.SshkeyaccessRemoved, _ = DiffGenericByHash(
		oldInv.SshKeyAccess,
		newInv.SshKeyAccess,
		func(s *schema.SshKeyAccess) string { return s.Fingerprint + s.AsUser },
	)

	delta.SshkeyinfoAdded, delta.SshkeyinfoRemoved, _ = DiffGenericByHash(
		oldInv.SshKeyInfo,
		newInv.SshKeyInfo,
		func(s *schema.SshKeyInfo) string { return s.Fingerprint },
	)

	return delta
}

// DiffGenericByHash computes the difference between two lists of proto messages based on their hashes.
func DiffGenericByHash[T proto.Message](
	oldList, newList []T,
	getKey func(T) string,
) (added, removed, changed []T) {
	oldHashes := make(map[string]uint64)
	oldMap := make(map[string]T)

	for _, o := range oldList {
		k := getKey(o)
		oldHashes[k] = hashMessage(o)
		oldMap[k] = o
	}

	for _, n := range newList {
		k := getKey(n)
		newHash := hashMessage(n)

		if oldHash, ok := oldHashes[k]; !ok {
			added = append(added, n)
		} else if oldHash != newHash {
			changed = append(changed, n)
		}
		delete(oldMap, k)
	}

	for _, o := range oldMap {
		removed = append(removed, o)
	}

	return
}

// IsDeltaEmpty checks if the delta inventory is empty
func IsDeltaEmpty(d *schema.HostDeltaInventory) bool {
	return len(d.PackagesAdded) == 0 &&
		len(d.PackagesRemoved) == 0 &&
		len(d.UsersAdded) == 0 &&
		len(d.UsersRemoved) == 0 &&
		len(d.SystemdservicesAdded) == 0 &&
		len(d.SystemdservicesRemoved) == 0 &&
		len(d.KnownhostsAdded) == 0 &&
		len(d.KnownhostsRemoved) == 0 &&
		len(d.SshkeyaccessAdded) == 0 &&
		len(d.SshkeyaccessRemoved) == 0 &&
		len(d.SshkeyinfoAdded) == 0 &&
		len(d.SshkeyinfoRemoved) == 0 &&
		len(d.SystemdservicesAdded) == 0 &&
		len(d.SystemdservicesRemoved) == 0 &&
		len(d.ProcessesAdded) == 0 &&
		len(d.ProcessesRemoved) == 0 &&
		d.Platform == nil &&
		d.Network == nil
}
