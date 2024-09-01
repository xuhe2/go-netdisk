package network

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
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

func (r *Repo) Init() error {
	// init repo
	if _, err := git.PlainInit(r.Path, false); err != nil {
		return err
	}
	return nil
}

func (r *Repo) AddRemote() error {
	repo, err := git.PlainOpen(r.Path)
	if err != nil {
		return err
	}
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{r.Url},
	})
	return err
}

func (r *Repo) Upload() error {
	// init repo if the repo is not initialized
	if _, err := git.PlainOpen(r.Path); err != nil {
		// init repo
		if err := r.Init(); err != nil {
			return err
		}
	}
	// add remote if the remote is not added
	if err := r.AddRemote(); err != nil {
		return err
	}
	return nil
}
