package network

import (
	"os"

	"github.com/go-git/go-git/v5"
)

type RepoConfig struct {
	Url  string
	Path string
}

type Repo struct {
	RepoConfig
}

func NewRepo(config RepoConfig) *Repo {
	return &Repo{config}
}

func (r *Repo) Download() error {
	_, err := git.PlainClone(r.Path, false, &git.CloneOptions{
		URL:      r.Url,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Upload() error {
	return nil
}
