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

func stableHashIgnoringUpdatedAt(msg proto.Message) uint64 {
	if msg == nil {
		return 0
	}

	cloned := proto.Clone(msg)
	resetUpdatedAtField(cloned)
	b, _ := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(cloned)

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

// DiffGenericByHash compares repeated protobuf entries by identity and content hash.
// Collection timestamps are metadata and do not make an item modified.
func DiffGenericByHash[T proto.Message](
	oldList, newList []T,
	getKey func(T) string,
) (added, removed, changed []T) {
	return diffGenericByHash(oldList, newList, getKey, stableHashIgnoringUpdatedAt)
}

func diffGenericByHash[T proto.Message](
	oldList, newList []T,
	getKey func(T) string,
	hashFn func(proto.Message) uint64,
) (added, removed, changed []T) {
	oldByKey := make(map[string]T, len(oldList))
	// Build a map of old items by their unique key.
	for _, item := range oldList {
		if isNilProtoValue(item) {
			continue
		}
		oldByKey[getKey(item)] = item
	}

	// Iterate over new items to find added and changed items.
	for _, item := range newList {
		if isNilProtoValue(item) {
			continue
		}

		// Check if the item exists in the old list.
		key := getKey(item)
		oldItem, exists := oldByKey[key]
		if !exists {
			added = append(added, item)
			continue
		}

		// If the item exists, check if it has changed by comparing stable hashes.
		if hashFn(oldItem) != hashFn(item) {
			changed = append(changed, item)
		}

		// Remove the item from the old map to track which items have been processed.
		delete(oldByKey, key)
	}

	// Any remaining items in oldByKey are considered removed.
	for _, item := range oldByKey {
		if isNilProtoValue(item) {
			continue
		}
		removed = append(removed, item)
	}

	return added, removed, changed
}

func uniqueByKey[T any](items []T, getKey func(T) string) []T {
	seen := make(map[string]struct{}, len(items))
	result := make([]T, 0, len(items))
	for _, item := range items {
		key := getKey(item)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, item)
	}
	return result
}

// Identity helpers for each entity.
func packageKey(item *schema.Package) string {
	if item == nil {
		return ""
	}
	return item.GetName() + ":" + item.GetArchitecture()
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

// packageVulnMatchKey returns a unique key for a PackageVulnMatch based on package name and installed version.
func packageVulnMatchKey(item *schema.PackageVulnMatch) string {
	if item == nil {
		return ""
	}
	return item.GetPackageName()
}

// computeComplianceReportDelta computes the delta between two compliance reports.
func computeComplianceReportDelta(oldReport, newReport *schema.ComplianceReport) *schema.ComplianceReportDelta {
	if oldReport == nil && newReport == nil {
		return nil
	}

	oldRules := []*schema.RuleCheckResult{}
	newRules := []*schema.RuleCheckResult{}

	if oldReport != nil {
		oldRules = oldReport.GetRuleResults()
	}
	if newReport != nil {
		newRules = newReport.GetRuleResults()
	}

	// Compute added and removed rules based on rule ID.
	rulesAdded, rulesRemoved, rulesChanged := DiffGenericByHash(
		oldRules,
		newRules,
		func(r *schema.RuleCheckResult) string { return r.GetId() },
	)

	delta := &schema.ComplianceReportDelta{
		RulesAdded:   rulesAdded,
		RulesRemoved: rulesRemoved,
		RulesChanged: rulesChanged,
	}

	if !sameProtoMessage(oldReport.GetScore(), newReport.GetScore()) {
		delta.ScoreDelta = newReport.GetScore()
	}
	if oldReport.GetProfile() != newReport.GetProfile() {
		delta.Profile = newReport.GetProfile()
	}
	if isComplianceReportDeltaEmpty(delta) {
		return nil
	}
	return delta
}

// computeVulnerabilityDelta computes the delta between two vulnerability reports.
func computeVulnerabilityDelta(oldReport, newReport *schema.VulnerabilityReport) *schema.VulnerabilityDelta {
	if oldReport == nil && newReport == nil {
		return nil
	}

	oldMatches := []*schema.PackageVulnMatch{}
	newMatches := []*schema.PackageVulnMatch{}

	if oldReport != nil {
		oldMatches = oldReport.GetMatches()
	}
	if newReport != nil {
		newMatches = newReport.GetMatches()
	}

	// Compute added and removed matches based on package name and installed version.
	matchesAdded, matchesRemoved, matchesChanged := DiffGenericByHash(
		oldMatches,
		newMatches,
		func(m *schema.PackageVulnMatch) string { return packageVulnMatchKey(m) },
	)

	if len(matchesAdded) == 0 && len(matchesRemoved) == 0 && len(matchesChanged) == 0 {
		return nil
	}

	return &schema.VulnerabilityDelta{
		MatchesAdded:   matchesAdded,
		MatchesRemoved: matchesRemoved,
		MatchesChanged: matchesChanged,
	}
}

// isVulnerabilityDeltaEmpty checks if a VulnerabilityDelta has no changes.
func isVulnerabilityDeltaEmpty(d *schema.VulnerabilityDelta) bool {
	if d == nil {
		return true
	}
	return len(d.GetMatchesAdded()) == 0 &&
		len(d.GetMatchesRemoved()) == 0 &&
		len(d.GetMatchesChanged()) == 0
}

// isComplianceReportDeltaEmpty checks if a ComplianceReportDelta has no changes.
func isComplianceReportDeltaEmpty(d *schema.ComplianceReportDelta) bool {
	if d == nil {
		return true
	}
	return len(d.GetRulesAdded()) == 0 &&
		len(d.GetRulesRemoved()) == 0 &&
		len(d.GetRulesChanged()) == 0 &&
		d.GetScoreDelta() == nil &&
		d.GetProfile() == ""
}

// Docker-specific key functions for granular diff.
func containerKey(item *schema.Containers) string {
	if item == nil {
		return ""
	}
	return item.GetId()
}

func imageKey(item *schema.ContainersImages) string {
	if item == nil {
		return ""
	}
	return item.GetId()
}

func dockerNetworkKey(item *schema.DockerNetworks) string {
	if item == nil {
		return ""
	}
	return item.GetId()
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
	delta.PackagesAdded, delta.PackagesRemoved, delta.PackagesChanged = DiffGenericByHash(
		oldInv.GetPackages(),
		newInv.GetPackages(),
		func(p *schema.Package) string { return packageKey(p) },
	)

	delta.UsersAdded, delta.UsersRemoved, delta.UsersChanged = DiffGenericByHash(
		oldInv.GetUsers(),
		newInv.GetUsers(),
		func(u *schema.User) string { return userKey(u) },
	)

	delta.ProcessesAdded, delta.ProcessesRemoved, delta.ProcessesChanged = DiffGenericByHash(
		oldInv.GetProcesses(),
		newInv.GetProcesses(),
		func(p *schema.Process) string { return processKey(p) },
	)

	delta.SystemdservicesAdded, delta.SystemdservicesRemoved, delta.SystemdservicesChanged = DiffGenericByHash(
		oldInv.GetSystemdService(),
		newInv.GetSystemdService(),
		func(s *schema.SystemdService) string { return systemdServiceKey(s) },
	)

	delta.KnownhostsAdded, delta.KnownhostsRemoved, delta.KnownhostsChanged = DiffGenericByHash(
		oldInv.GetKnownHost(),
		newInv.GetKnownHost(),
		func(k *schema.KnownHost) string { return knownHostKey(k) },
	)

	delta.SshkeyinfoAdded, delta.SshkeyinfoRemoved, delta.SshkeyinfoChanged = diffGenericByHash(
		uniqueByKey(oldInv.GetSshKeyInfo(), sshKeyInfoKey),
		uniqueByKey(newInv.GetSshKeyInfo(), sshKeyInfoKey),
		func(k *schema.SshKeyInfo) string { return sshKeyInfoKey(k) },
		stableHashIgnoringUpdatedAt,
	)

	delta.SshkeyaccessAdded, delta.SshkeyaccessRemoved, _ = diffGenericByHash(
		uniqueByKey(oldInv.GetSshKeyAccess(), sshKeyAccessKey),
		uniqueByKey(newInv.GetSshKeyAccess(), sshKeyAccessKey),
		func(k *schema.SshKeyAccess) string { return sshKeyAccessKey(k) },
		stableHashIgnoringUpdatedAt,
	)

	// Compute Docker delta with granular container/image/network changes.
	dockerDelta := computeDockerDelta(oldInv.GetApplication(), newInv.GetApplication())
	if dockerDelta != nil {
		delta.DockerDelta = dockerDelta
	}

	vulnDelta := computeVulnerabilityDelta(oldInv.GetVulnerabilityReport(), newInv.GetVulnerabilityReport())
	if vulnDelta != nil {
		delta.VulnerabilityDelta = vulnDelta
	}

	complianceDelta := computeComplianceReportDelta(oldInv.GetComplianceReport(), newInv.GetComplianceReport())
	if complianceDelta != nil {
		delta.ComplianceReportDelta = complianceDelta
	}

	return delta
}

// IsDeltaEmpty returns true when no meaningful change is present.
func IsDeltaEmpty(d *schema.HostDeltaInventory) bool {
	if d == nil {
		return true
	}

	return len(d.GetPackagesAdded()) == 0 &&
		len(d.GetPackagesRemoved()) == 0 &&
		len(d.GetPackagesChanged()) == 0 &&
		len(d.GetUsersAdded()) == 0 &&
		len(d.GetUsersChanged()) == 0 &&
		len(d.GetUsersRemoved()) == 0 &&
		isDockerDeltaEmpty(d.GetDockerDelta()) &&
		len(d.GetSystemdservicesAdded()) == 0 &&
		len(d.GetSystemdservicesRemoved()) == 0 &&
		len(d.GetSystemdservicesChanged()) == 0 &&
		len(d.GetKnownhostsAdded()) == 0 &&
		len(d.GetKnownhostsRemoved()) == 0 &&
		len(d.GetKnownhostsChanged()) == 0 &&
		len(d.GetSshkeyaccessAdded()) == 0 &&
		len(d.GetSshkeyaccessRemoved()) == 0 &&
		len(d.GetSshkeyinfoAdded()) == 0 &&
		len(d.GetSshkeyinfoRemoved()) == 0 &&
		len(d.GetSshkeyinfoChanged()) == 0 &&
		len(d.GetProcessesAdded()) == 0 &&
		len(d.GetProcessesRemoved()) == 0 &&
		len(d.GetProcessesChanged()) == 0 &&
		isVulnerabilityDeltaEmpty(d.GetVulnerabilityDelta()) &&
		isComplianceReportDeltaEmpty(d.GetComplianceReportDelta()) &&
		d.GetPlatform() == nil &&
		d.GetNetwork() == nil
}

// computeDockerDelta computes granular Docker delta from application lists.
func computeDockerDelta(oldApps, newApps []*schema.Application) *schema.DockerDelta {
	oldDocker := extractDocker(oldApps)
	newDocker := extractDocker(newApps)

	if oldDocker == nil && newDocker == nil {
		return nil
	}

	if oldDocker == nil {
		oldDocker = &schema.Docker{}
	}
	if newDocker == nil {
		newDocker = &schema.Docker{}
	}

	containersAdded, containersRemoved, containersChanged := DiffGenericByHash(
		oldDocker.GetContainers(),
		newDocker.GetContainers(),
		func(c *schema.Containers) string { return containerKey(c) },
	)

	imagesAdded, imagesRemoved, imagesChanged := DiffGenericByHash(
		oldDocker.GetImages(),
		newDocker.GetImages(),
		func(i *schema.ContainersImages) string { return imageKey(i) },
	)

	networksAdded, networksRemoved, networksChanged := DiffGenericByHash(
		oldDocker.GetNetworks(),
		newDocker.GetNetworks(),
		func(n *schema.DockerNetworks) string { return dockerNetworkKey(n) },
	)

	if len(containersAdded) == 0 && len(containersRemoved) == 0 && len(containersChanged) == 0 &&
		len(imagesAdded) == 0 && len(imagesRemoved) == 0 && len(imagesChanged) == 0 &&
		len(networksAdded) == 0 && len(networksRemoved) == 0 && len(networksChanged) == 0 {
		return nil
	}

	return &schema.DockerDelta{
		ContainersAdded:   containersAdded,
		ContainersRemoved: containersRemoved,
		ContainersChanged: containersChanged,
		ImagesAdded:       imagesAdded,
		ImagesRemoved:     imagesRemoved,
		ImagesChanged:     imagesChanged,
		NetworksAdded:     networksAdded,
		NetworksRemoved:   networksRemoved,
		NetworksChanged:   networksChanged,
	}
}

// extractDocker extracts the Docker struct from an Application list (assumes only one Docker app).
func extractDocker(apps []*schema.Application) *schema.Docker {
	if len(apps) == 0 {
		return nil
	}
	for _, app := range apps {
		if app != nil && app.GetDocker() != nil {
			return app.GetDocker()
		}
	}
	return nil
}

// isDockerDeltaEmpty returns true when the Docker delta has no changes.
func isDockerDeltaEmpty(d *schema.DockerDelta) bool {
	if d == nil {
		return true
	}
	return len(d.GetContainersAdded()) == 0 &&
		len(d.GetContainersRemoved()) == 0 &&
		len(d.GetContainersChanged()) == 0 &&
		len(d.GetImagesAdded()) == 0 &&
		len(d.GetImagesRemoved()) == 0 &&
		len(d.GetImagesChanged()) == 0 &&
		len(d.GetNetworksAdded()) == 0 &&
		len(d.GetNetworksRemoved()) == 0 &&
		len(d.GetNetworksChanged()) == 0
}
