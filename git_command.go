package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Git is a struct
type Git struct {
	dir string
}

func newGit(dir string) (*Git, error) {
	if !isDir(dir) {
		return nil, errors.New("invalid parameter")
	}
	git := &Git{dir}

	if git.isInsideWorkTree() {
		topDir, _ := git.getTopDir()
		git.dir = topDir
	}
	return git, nil
}

// git rev-parse --show-toplevel
func (git Git) getTopDir() (string, error) {
	return git.exec("rev-parse", "--show-toplevel")
}

// git config --get remote.origin.url
//   => git@github.com:inouet/gh-open.git
func (git Git) getRemoteOriginURL() (string, error) {
	return git.exec("config", "--get", "remote.origin.url")
}

// git rev-parse HEAD
//   => 695895662d96bac8d94fd71dc9d2dec534c8e494
func (git Git) getCommitHash() (string, error) {
	return git.exec("rev-parse", "HEAD")
}

// git rev-parse --is-inside-work-tree
//   => true or false
func (git Git) isInsideWorkTree() bool {
	boolStr, err := git.exec("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false
	}
	if boolStr == "true" {
		return true
	}
	return false
}

// git clone
func (git Git) clone(repo string) (string, error) {
	return git.exec("clone", repo)
}

// git config name
func (git Git) getConfig(name, defaultValue string) string {
	configValue, err := git.exec("config", "--get", name)
	if err != nil {
		return defaultValue
	}
	return configValue
}

func (git Git) exec(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = git.dir
	out, err := cmd.Output()

	if err != nil {
		argStr := strings.Join(args, " ")
		return "", fmt.Errorf("git command failed: git %s", argStr)
	}

	return strings.TrimSpace(string(out)), nil
}
