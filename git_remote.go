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

	vcsurl "github.com/gitsight/go-vcsurl"
)

const (
	gitConfigURLTypeName  string = "gh-open.urltype"
	gitConfigProtocolName string = "gh-open.protocol"
)

// GitRemote is a struct
type GitRemote struct {
	git  *Git
	path string // Relative path from git top directory
}

// Enhanced regex to match both standard git SSH URLs and organization format URLs
var orgSSHRegex = regexp.MustCompile(`^(git|org-[a-zA-Z0-9_-]+)@([a-zA-Z0-9._-]+):([a-zA-Z0-9/_-]+)(/[a-zA-Z0-9/_-]+)*(.git)?$`)

// parseSSHRemoteURL parses SSH remote URLs including both standard format and organization format
// e.g. git@github.com:user/repo.git or org-1234@github.com:user/repo.git
func parseSSHRemoteURL(remoteURL string) (host, username, repo string, err error) {
	matches := orgSSHRegex.FindStringSubmatch(remoteURL)
	if matches == nil {
		return "", "", "", fmt.Errorf("invalid SSH remote URL format: %s", remoteURL)
	}

	host = matches[2]
	repoPath := matches[3]

	// Remove .git suffix if present
	repoPath = strings.TrimSuffix(repoPath, ".git")

	// Split the path into username and repo name
	parts := strings.Split(repoPath, "/")
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid repository path: %s", repoPath)
	}

	username = parts[0]
	repo = parts[len(parts)-1]

	return host, username, repo, nil
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
	urlType := r.git.getConfig(gitConfigURLTypeName, "")
	scheme := r.git.getConfig(gitConfigProtocolName, "https")

	// Try to handle special SSH URL format with organization ID (org-ID@github.com:user/repo.git)
	if strings.Contains(remote, "@") && strings.Contains(remote, ":") && !strings.HasPrefix(remote, "git@") {
		host, username, repo, err := parseSSHRemoteURL(remote)
		if err == nil {
			// Create URL object
			newURL := neturl.URL{
				Scheme: scheme,
				Host:   host,
				Path:   "/" + username + "/" + repo,
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
		// If parsing fails, fall back to the standard method
	}

	// Parse with gitsight/go-vcsurl for standard formats
	info, err := vcsurl.Parse(remote)
	if err != nil {
		return "", err
	}

	// Create new URL with the parsed information
	newURL := neturl.URL{
		Scheme: scheme,
		Host:   string(info.Host),
		Path:   "/" + info.FullName,
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
