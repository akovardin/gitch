package syncer

import (
	"errors"
	"fmt"
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
	publicKeys, err := ssh.NewPublicKeys("git", []byte(s.key), "")
	if err != nil {
		return fmt.Errorf("error on create key: %w", err)
	}

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		Auth:     publicKeys,
		URL:      s.from,
		Progress: os.Stdout,
		Mirror:   true,
		Tags:     git.AllTags,
	})

	if err != nil {
		return fmt.Errorf("error on clone: %w", err)
	}

	remote, err := r.CreateRemote(&config.RemoteConfig{
		Name:   "sync",
		URLs:   []string{s.to},
		Mirror: true,
	})
	if err != nil {
		return fmt.Errorf("error on create remote: %w", err)
	}

	if err := remote.Push(&git.PushOptions{
		FollowTags: true,
		Auth:       publicKeys,
		RemoteName: "sync",
		Force:      true,
		Progress:   os.Stdout,
	}); err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("error on push: %w", err)
	}

	return nil
}
