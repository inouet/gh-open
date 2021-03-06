package main

import (
	"errors"
	"fmt"
	neturl "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	giturls "github.com/whilp/git-urls"
)

const (
	gitConfigUrlTypeName  string = "gh-open.urltype"
	gitConfigProtocolName string = "gh-open.protocol"
)

// GitRemote is a struct
type GitRemote struct {
	git  *Git
	path string // Relative path from git top directory
}

func newGitRemote(objectPath string) (*GitRemote, error) {
	isFile := isFile(objectPath)
	isDir := isDir(objectPath)

	if !isFile && !isDir {
		return nil, fmt.Errorf("%s: no such file or directory", objectPath)
	}

	absPath, _ := filepath.Abs(objectPath)
	absDir := absPath
	if isFile {
		absDir = filepath.Dir(absPath)
	}

	g, err := newGit(absDir)

	if err != nil {
		return nil, err
	}

	if !g.isInsideWorkTree() {
		return nil, errors.New("not a git repository (or any of the parent directories)")
	}

	relPath, err := relativePath(g.dir, absPath)
	if err != nil {
		return nil, err
	}
	gitRemote := GitRemote{g, relPath}

	return &gitRemote, nil
}

func (r GitRemote) remoteURL(branch string, line1, line2 int) (string, error) {
	remote, err := r.git.getRemoteOriginURL()
	if err != nil {
		return remote, err
	}

	// If it cannot be determined from the remote domain,
	//   read the setting from git config and make a judgment based on it.
	urlType := r.git.getConfig(gitConfigUrlTypeName, "")
	scheme := r.git.getConfig(gitConfigProtocolName, "https")

	url, err := giturls.Parse(remote)
	if err != nil {
		return "", err
	}

	// remove .git from /inouet/gh-open.git
	path := strings.TrimSuffix(url.Path, ".git")
	newURL := neturl.URL{
		Scheme: scheme,
		Host:   url.Host,
		Path:   path,
	}

	if r.path == "" && branch == "" {
		return newURL.String(), nil
	}

	if branch == "" {
		branch, err = r.git.getCommitHash()
		if err != nil {
			return "", err
		}
	}

	remoteURL, err := buildURL(newURL, r.path, branch, line1, line2, urlType)
	if err != nil {
		return "", err
	}

	return remoteURL, nil
}

func isFile(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	if fi.Mode().IsRegular() {
		return true
	}
	return false
}

func isDir(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	if fi.Mode().IsDir() {
		return true
	}
	return false
}

func relativePath(obj1, obj2 string) (string, error) {
	o1, _ := filepath.EvalSymlinks(obj1)
	o2, _ := filepath.EvalSymlinks(obj2)
	relPath, err := filepath.Rel(o1, o2)
	if err != nil {
		return "", err
	}

	if relPath == "." {
		relPath = ""
	}
	return relPath, nil
}

// valid format: 20 or 20-30
func getLineOption(line string) (int, int, error) {
	line = strings.TrimSpace(line)
	line1, line2 := 0, 0
	if line == "" {
		return line1, line2, nil
	}
	matched, err := regexp.MatchString(`^([0-9]+|[0-9]+-[0-9]+)$`, line)
	if err != nil || !matched {
		return line1, line2, errors.New("invalid line format")
	}
	arr := strings.Split(line, "-")
	if len(arr) == 1 {
		line1, _ = strconv.Atoi(arr[0])
	}
	if len(arr) == 2 {
		line1, _ = strconv.Atoi(arr[0])
		line2, _ = strconv.Atoi(arr[1])
	}
	return line1, line2, nil
}
