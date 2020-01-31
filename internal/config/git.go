/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package config

import (
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
func SyncGitConfigs(customDir string) bool {
	if load.Args.GitToken != "" && load.Args.GitUser != "" && load.Args.GitRepo != "" {
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
			load.Logrus.WithFields(logrus.Fields{
				"repo": load.Args.GitRepo,
				"err":  err,
			}).Error("config: git sync invalid url")
		} else {
			repoDir := path.Join(syncDir, u.Path)
			_, err = appFS.Stat(repoDir)
			if err == nil {
				load.Logrus.WithFields(logrus.Fields{
					"repo": load.Args.GitRepo,
				}).Debug("config: git sync pulling repo")

				err := GitPull(repoDir)
				if err != nil {
					load.Logrus.WithFields(logrus.Fields{
						"err":  err,
						"repo": load.Args.GitRepo,
					}).Error("config: git sync pull failed")
				}
				if err == nil {
					return true
				}
			} else {
				load.Logrus.WithFields(logrus.Fields{
					"repo": load.Args.GitRepo,
				}).Debug("config: git sync cloning repo")

				err := GitClone(repoDir, u)
				if err != nil {
					load.Logrus.WithFields(logrus.Fields{
						"err":  err,
						"repo": load.Args.GitRepo,
					}).Error("config: git clone failed")
				}

				if err == nil {
					return true
				}
			}
		}
	} else {
		load.Logrus.Debug("config: git sync configuration not set")
	}

	return false
}

// GitClone git clone
func GitClone(dir string, u *url.URL) error {
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: u.Scheme + "://" + load.Args.GitUser + ":" + load.Args.GitToken + "@" + u.Host + u.Path,
		// Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	GitCheckout(w)

	return nil
}

// GitPull git pull
func GitPull(dir string) error {
	// instance\iate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	GitCheckout(w)

	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		if err.Error() == "already up-to-date" {
			load.Logrus.Debug("config: git pull - " + dir + " " + err.Error())
			return nil
		}
		return err
	}

	return nil
}

// GitCheckout git checkout
func GitCheckout(w *git.Worktree) {
	if load.Args.GitBranch != "" && load.Args.GitCommit == "" {
		err := w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(load.Args.GitBranch),
		})
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"err":    err,
				"repo":   load.Args.GitRepo,
				"branch": load.Args.GitBranch,
			}).Error("config: git sync checkout failed")
		}
	} else if load.Args.GitCommit != "" {
		err := w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(load.Args.GitCommit),
		})
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"err":    err,
				"repo":   load.Args.GitRepo,
				"commit": load.Args.GitCommit,
			}).Error("config: git sync checkout failed")
		}
	}
}
