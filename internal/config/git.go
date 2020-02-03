/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var appFS = afero.NewOsFs()

// SyncGitConfigs Clone git repo if already exists, else pull latest version
func SyncGitConfigs(customDir string) (bool, error) {
	if load.Args.GitToken == "" ||
		load.Args.GitUser == "" ||
		load.Args.GitRepo == "" {
		load.Logrus.Debug("config: git sync configuration not set")

		return false, nil
	}
	syncDir := load.Args.ConfigDir
	if customDir != "" {
		syncDir = customDir
	}

	if !strings.HasSuffix(load.Args.GitRepo, "/") {
		load.Args.GitRepo = load.Args.GitRepo + "/"
	}

	load.Logrus.Debugf("config: syncing git configs %v %v into %v", load.Args.GitService, load.Args.GitRepo, syncDir)

	u, err := url.Parse(load.Args.GitRepo)
	if err != nil {
		return false, fmt.Errorf("config: git sync invalid url, repo: %s, error: %v ", load.Args.GitRepo, err)
	}
	repoDir := path.Join(syncDir, u.Path)
	_, err = appFS.Stat(repoDir)

	// If cannot access the repo dir, clone it.
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"repo": load.Args.GitRepo,
		}).Debug("config: git sync cloning repo")

		err := GitClone(repoDir, u)
		if err != nil {
			return false, fmt.Errorf("config: git clone failed, repo: %s, error: %v ", load.Args.GitRepo, err)
		}
	} else {
		load.Logrus.WithFields(logrus.Fields{
			"repo": load.Args.GitRepo,
		}).Debug("config: git sync pulling repo")

		err = GitPull(repoDir)
		if err != nil {
			return false, fmt.Errorf("config: git sync pull failed, repo: %s, error: %v ", load.Args.GitRepo, err)
		}
	}
	return true, nil
}

// GitClone git clone
func GitClone(dir string, u *url.URL) error {
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: u.Scheme + "://" + load.Args.GitUser + ":" + load.Args.GitToken + "@" + u.Host + u.Path,
		// Progress: os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repo, error: %v", err)
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get working directory while cloning the repo, error: %v", err)
	}

	err = GitCheckout(w)
	if err != nil {
		return fmt.Errorf("failed to clone repo, error: %v", err)
	}

	return nil
}

// GitPull git pull
func GitPull(dir string) error {
	// instance\iate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("failed to pull from repo, error: %v", err)
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get working directory while pulling the repo, error: %v", err)
	}

	err = GitCheckout(w)
	if err != nil {
		return fmt.Errorf("failed to pull from repo, error: %v", err)
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		if err.Error() == "already up-to-date" {
			load.Logrus.Debug("config: git pull - " + dir + " " + err.Error())
			return nil
		}

		return fmt.Errorf("error occurred while pulling the repo, error: %v", err)
	}

	return nil
}

// GitCheckout git checkout
func GitCheckout(w *git.Worktree) error {
	if load.Args.GitBranch != "" && load.Args.GitCommit == "" {
		err := w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(load.Args.GitBranch),
		})
		if err != nil {
			return fmt.Errorf("failed to checkout repo: %s, branch: %s, error: %v", load.Args.GitRepo, load.Args.GitBranch, err)
		}
	} else if load.Args.GitCommit != "" {
		err := w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(load.Args.GitCommit),
		})
		if err != nil {
			return fmt.Errorf("failed to checkout repo: %s, commit: %s, error: %v", load.Args.GitRepo, load.Args.GitCommit, err)
		}
	}
	return nil
}
