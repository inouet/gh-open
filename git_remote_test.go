package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func mkTempDir() string {
	tmpDir := os.TempDir()
	testDir, err := ioutil.TempDir(tmpDir, "tests-*")
	if err != nil {
		log.Fatal(err)
	}
	return testDir
}

func TestRemoteUrl(t *testing.T) {

	cases := map[string]struct {
		path   string
		branch string
		line   string
		want   string
	}{
		"root-dir":
		{
			path:   "./",
			branch: "",
			line:   "",
			want:   "https://github.com/inouet/gh-open",
		},
		"root-dir-master":
		{
			path:   "./",
			branch: "master",
			line:   "",
			want:   "https://github.com/inouet/gh-open/tree/master/",
		},
		"readme.md":
		{
			path:   "./README.md",
			branch: "master",
			line:   "10",
			want:   "https://github.com/inouet/gh-open/tree/master/README.md#L10",
		},
	}

	for name, c := range cases {
		gr, err := newGitRemote(c.path)
		if err != nil {
			t.Fatal()
		}
		got, _ := gr.remoteURL(c.branch, c.line)
		if got != c.want {
			t.Errorf("%s want '%s', got '%s'\n", name, c.want, got)
		}
	}
}

func TestRemoteUrlFunctional(t *testing.T) {

	t.Parallel()
	cloneDir := mkTempDir()
	defer os.RemoveAll(cloneDir)

	cases := map[string]struct {
		repo    string
		repoDir string
		path    string
		branch  string
		line    string
		want    string
	}{
		"github-cheat-sheet":
		{
			repo:    "https://github.com/githubtraining/github-cheat-sheet.git",
			repoDir: "github-cheat-sheet",
			path:    "LICENSE",
			branch:  "master",
			line:    "3-4",
			want:    "https://github.com/githubtraining/github-cheat-sheet/tree/master/LICENSE#L3-4",
		},
		"bitbucket-test":
		{
			repo:    "https://bitbucket.org/atn13/bitbucketstationlocations.git",
			repoDir: "bitbucketstationlocations",
			path:    "README.txt",
			branch:  "master",
			line:    "2-4",
			want:    "https://bitbucket.org/atn13/bitbucketstationlocations/src/master/README.txt#lines-2:4",
		},
		"gitlab-gitlab-examples-docker":
		{
			repo:    "https://gitlab.com/gitlab-examples/docker.git",
			repoDir: "docker",
			path:    "Dockerfile",
			branch:  "master",
			line:    "1",
			want:    "https://gitlab.com/gitlab-examples/docker/-/blob/master/Dockerfile#L1",
		},
	}

	git, _ := newGit(cloneDir)

	for name, c := range cases {
		c := c
		git.clone(c.repo)
		gitDir := filepath.Join(cloneDir, c.repoDir)
		path := filepath.Join(gitDir, c.path)
		gr, err := newGitRemote(path)
		if err != nil {
			t.Fatal()
		}
		got, err := gr.remoteURL(c.branch, c.line)
		if err != nil {
			t.Fatal()
		}
		if got != c.want {
			t.Errorf("%s want '%s', got '%s'\n", name, c.want, got)
		}
	}
}

func TestNewGitRemote(t *testing.T) {

	emptyDir := mkTempDir()
	defer os.RemoveAll(emptyDir)

	cases := []struct {
		path    string
		want    string
		wantErr error
	}{
		{
			path:    emptyDir,
			wantErr: errors.New("not a git repository (or any of the parent directories)"),
		},
	}

	for _, c := range cases {
		_, err := newGitRemote(c.path)
		if c.wantErr != nil {
			if err == nil {
				t.Errorf("want '%s', got '%s'\n", c.wantErr.Error(), err)
				continue
			}
			if c.wantErr.Error() != err.Error() {
				t.Errorf("want '%s', got '%s'\n", c.wantErr.Error(), err.Error())
			}
		}
	}
}

func TestRelativePath(t *testing.T) {
	testDir := mkTempDir()
	defer os.RemoveAll(testDir)

	subDir := "foo/bar/baz"
	os.MkdirAll(testDir+"/"+subDir, 0777)

	cases := []struct {
		path1 string
		path2 string
		want  string
	}{
		{
			path1: testDir,
			path2: testDir + "/" + subDir,
			want:  subDir,
		},
		{
			path1: testDir,
			path2: testDir,
			want:  "", // same directory return ""
		},
	}
	for _, c := range cases {
		got, _ := relativePath(c.path1, c.path2)
		if got != c.want {
			t.Errorf("want '%s', got '%s'\n", c.want, got)
		}
	}
}

func TestValidateLine(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{line: "", want: true},
		{line: "3", want: true},
		{line: "3-10", want: true},
		{line: "3-", want: false},
	}
	for _, c := range cases {
		got := validateLine(c.line)
		if got != c.want {
			t.Errorf("'%s' want %v, got %v\n", c.line, c.want, got)
		}
	}
}