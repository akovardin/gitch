package syncer

import (
	"errors"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Syncer struct {
	key  string
	from string
	to   string
}

func New(from, to, key string) *Syncer {
	return &Syncer{
		key:  key,
		from: from,
		to:   to,
	}
}

func (s *Syncer) Sync() error {
	sshKey, err := os.ReadFile(s.key)
	if err != nil {
		return err
	}

	publicKeys, err := ssh.NewPublicKeys("git", sshKey, "")
	if err != nil {
		return err
	}

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		Auth:     publicKeys,
		URL:      s.from,
		Progress: os.Stdout,
		Mirror:   true,
	})

	if err != nil {
		return err
	}

	remote, err := r.CreateRemote(&config.RemoteConfig{
		Name:   "sync",
		URLs:   []string{s.to},
		Mirror: true,
	})
	if err != nil {
		return err
	}

	if err := remote.Push(&git.PushOptions{
		Auth:       publicKeys,
		RemoteName: "sync",
		Force:      true,
		Progress:   os.Stdout,
	}); err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil
}
