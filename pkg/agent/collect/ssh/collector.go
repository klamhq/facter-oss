package ssh

import (
	"context"

	"github.com/klamhq/facter-oss/pkg/agent/collectors/ssh"
	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
)

type SSHInfosCollectorImpl struct {
	log *logrus.Logger
	cfg *options.SSHOptions
}

func New(log *logrus.Logger, cfg *options.SSHOptions) *SSHInfosCollectorImpl {

	return &SSHInfosCollectorImpl{
		log: log,
		cfg: cfg,
	}
}

func (c *SSHInfosCollectorImpl) Name() string { return "ssh" }

func (c *SSHInfosCollectorImpl) CollectSSHInfos(ctx context.Context, users []*schema.User) ([]*schema.SshKeyAccess, []*schema.KnownHost, []*schema.SshKeyInfo, error) {
	c.log.Info("Crafting sshkeys and known hosts")
	var homeDirList []string

	for _, u := range users {
		homeDirList = append(homeDirList, u.Home)
	}
	sshKeysPubInfo, sshAuthorizedKeyInfo, knownHost := ssh.GetSshInfo(c.log, homeDirList)
	collectedSshKeysInfos := make([]*schema.SshKeyInfo, 0, len(sshKeysPubInfo))
	for _, sshKey := range sshKeysPubInfo {
		protoSsh := &schema.SshKeyInfo{
			Fingerprint:        sshKey.Fingerprint,
			Type:               sshKey.Type,
			Length:             sshKey.Length,
			Comment:            sshKey.Comment,
			Path:               sshKey.Path,
			Name:               sshKey.Name,
			FromAuthorizedKeys: sshKey.FromAuthorizedKeys,
			Owner:              sshKey.Owner,
		}

		for _, opts := range sshKey.Options {
			protoSsh.Options = append(protoSsh.Options, &schema.SshKeyOptions{
				Options: opts,
			})
		}
		collectedSshKeysInfos = append(collectedSshKeysInfos, protoSsh)
	}
	collectedSshKeysAccess := make([]*schema.SshKeyAccess, 0, len(sshAuthorizedKeyInfo))
	for _, sshKeyAccess := range sshAuthorizedKeyInfo {
		protoSsh := &schema.SshKeyAccess{
			Fingerprint: sshKeyAccess.Fingerprint,
			AsUser:      sshKeyAccess.AsUser,
		}

		collectedSshKeysAccess = append(collectedSshKeysAccess, protoSsh)
	}
	collectedKnownHost := make([]*schema.KnownHost, 0, len(knownHost))
	for _, entry := range knownHost {
		protoKnownHost := &schema.KnownHost{
			Hostname:    entry.Hostname,
			Type:        entry.Type,
			Fingerprint: entry.Fingerprint,
			Owner:       entry.Owner,
		}
		collectedKnownHost = append(collectedKnownHost, protoKnownHost)
	}
	return collectedSshKeysAccess, collectedKnownHost, collectedSshKeysInfos, nil
}
