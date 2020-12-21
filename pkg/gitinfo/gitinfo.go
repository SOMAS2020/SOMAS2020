<<<<<<< HEAD
// Package gitinfo gets information about the git repo.
package gitinfo

import (
	"fmt"
	"strings"

	"github.com/SOMAS2020/SOMAS2020/pkg/sysutils"
	"github.com/pkg/errors"
)

// GitInfo contains information about the latest config of the repo.
type GitInfo struct {
	Hash      string
	ShortHash string
	GithubURL string
}

// GetGitInfo gets GitInfo in the specified working directory wd.
func GetGitInfo(wd string) (GitInfo, error) {
	gitInfo := GitInfo{}

	hashBuf, err := sysutils.RunCommandInDir("git", []string{"rev-parse", "HEAD"}, wd)
	if err != nil {
		return gitInfo, errors.Errorf("Failed to get git hash: %v", err)
	}
	hash := strings.TrimSpace(string(hashBuf))
	gitInfo.Hash = hash
	gitInfo.ShortHash = hash[:7]

	gitInfo.GithubURL = fmt.Sprintf("https://github.com/SOMAS2020/SOMAS2020/tree/%v", hash)

	return gitInfo, nil
}
=======
// Package gitinfo gets information about the git repo.
package gitinfo

import (
	"fmt"
	"strings"

	"github.com/SOMAS2020/SOMAS2020/pkg/sysutils"
	"github.com/pkg/errors"
)

// GitInfo contains information about the latest config of the repo.
type GitInfo struct {
	Hash      string
	ShortHash string
	GithubURL string
}

// GetGitInfo gets GitInfo in the specified working directory wd.
func GetGitInfo(wd string) (GitInfo, error) {
	gitInfo := GitInfo{}

	hashBuf, err := sysutils.RunCommandInDir("git", []string{"rev-parse", "HEAD"}, wd)
	if err != nil {
		return gitInfo, errors.Errorf("Failed to get git hash: %v", err)
	}
	hash := strings.TrimSpace(string(hashBuf))
	gitInfo.Hash = hash
	gitInfo.ShortHash = hash[:7]

	remoteURLBuf, err := sysutils.RunCommandInDir("git", []string{"config", "--get", "remote.origin.url"}, wd)
	if err != nil {
		return gitInfo, errors.Errorf("Failed to get git remote origin url: %v", err)
	}
	remoteURL := strings.TrimSpace(string(remoteURLBuf))

	gitInfo.GithubURL = fmt.Sprintf("%v/tree/%v", string(remoteURL), hash)

	return gitInfo, nil
}
>>>>>>> main
