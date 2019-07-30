package config

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var appFS = afero.NewOsFs()

// SyncGitConfigs Clone git repo if already exists, else pull latest version
func SyncGitConfigs(customDir string) bool {
	if load.Args.GitService != "" && load.Args.GitToken != "" && load.Args.GitUser != "" && load.Args.GitRepo != "" {
		syncDir := load.Args.ConfigDir
		if customDir != "" {
			syncDir = customDir
		}

		if !strings.HasSuffix(load.Args.GitRepo, "/") {
			load.Args.GitRepo = load.Args.GitRepo + "/"
		}

		logger.Flex("debug", nil, fmt.Sprintf("syncing git configs %v %v into %v", load.Args.GitService, load.Args.GitRepo, syncDir), false)

		u, err := url.Parse(load.Args.GitRepo)
		if err != nil {
			logger.Flex("error", err, "invalid url", false)
		} else {
			repoDir := path.Join(syncDir, u.Path)
			_, err = appFS.Stat(repoDir)
			if err == nil {
				logger.Flex("debug", nil, fmt.Sprintf("pulling git repo %v", load.Args.GitRepo), false)
				err := GitPull(repoDir)
				logger.Flex("error", err, "git pull failed", false)
				if err == nil {
					return true
				}
			} else {
				logger.Flex("debug", nil, fmt.Sprintf("cloning git repo %v", load.Args.GitRepo), false)
				err := GitClone(repoDir, u)
				logger.Flex("error", err, "", false)
				if err == nil {
					return true
				}
			}
		}
	} else {
		logger.Flex("debug", nil, "git sync configuration not set", false)
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
			logger.Flex("debug", nil, "git pull - "+dir+" "+err.Error(), false)
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
		logger.Flex("error", err, load.Args.GitBranch, false)
	} else if load.Args.GitCommit != "" {
		err := w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(load.Args.GitCommit),
		})
		logger.Flex("error", err, load.Args.GitCommit, false)
	}
}
