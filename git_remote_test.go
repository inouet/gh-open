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
		line1  int
		line2  int
		want   string
	}{
		"root-dir":
		{
			path:   "./",
			branch: "",
			line1:  0,
			line2:  0,
			want:   "https://github.com/inouet/gh-open",
		},
		"root-dir-master":
		{
			path:   "./",
			branch: "master",
			line1:  0,
			line2:  0,
			want:   "https://github.com/inouet/gh-open/tree/master/",
		},
		"readme.md":
		{
			path:   "./README.md",
			branch: "master",
			line1:  10,
			line2:  0,
			want:   "https://github.com/inouet/gh-open/tree/master/README.md#L10",
		},
	}

	for name, c := range cases {
		gr, err := newGitRemote(c.path)
		if err != nil {
			t.Fatal(err)
		}
		got, _ := gr.remoteURL(c.branch, c.line1, c.line2)
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
		line1   int
		line2   int
		want    string
	}{
		"github-cheat-sheet":
		{
			repo:    "https://github.com/githubtraining/github-cheat-sheet.git",
			repoDir: "github-cheat-sheet",
			path:    "LICENSE",
			branch:  "master",
			line1:   3,
			line2:   4,
			want:    "https://github.com/githubtraining/github-cheat-sheet/tree/master/LICENSE#L3-L4",
		},
		"bitbucket-test":
		{
			repo:    "https://bitbucket.org/atn13/bitbucketstationlocations.git",
			repoDir: "bitbucketstationlocations",
			path:    "README.txt",
			branch:  "master",
			line1:   2,
			line2:   4,
			want:    "https://bitbucket.org/atn13/bitbucketstationlocations/src/master/README.txt#lines-2:4",
		},
		"gitlab-gitlab-examples-docker":
		{
			repo:    "https://gitlab.com/gitlab-examples/docker.git",
			repoDir: "docker",
			path:    "Dockerfile",
			branch:  "master",
			line1:   1,
			line2:   2,
			want:    "https://gitlab.com/gitlab-examples/docker/-/blob/master/Dockerfile#L1-2",
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
			t.Fatal(err)
		}
		got, err := gr.remoteURL(c.branch, c.line1, c.line2)
		if err != nil {
			t.Fatal(err)
		}
		if got != c.want {
			t.Errorf("%s want '%s', got '%s'\n", name, c.want, got)
		}
	}
}

func TestConfig(t *testing.T) {
	testDir := mkTempDir()
	defer os.RemoveAll(testDir)

	git, _ := newGit(testDir)
	git.clone("https://github.com/githubtraining/github-cheat-sheet.git")

	path := filepath.Join(testDir, "/github-cheat-sheet")
	git, _ = newGit(path)

	// set config
	git.exec("config", gitConfigUrlTypeName, "bitbucket.org")
	git.exec("config", gitConfigProtocolName, "http")

	gr, err := newGitRemote(filepath.Join(path, "LICENSE"))
	if err != nil {
		t.Fatal(err)
	}
	got, _ := gr.remoteURL("master", 3, 4)

	// Expect bitbucket style url and http protocol
	want := "http://github.com/githubtraining/github-cheat-sheet/src/master/LICENSE#lines-3:4"
	if got != want {
		t.Errorf("want '%s', got '%s'\n", want, got)
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

func TestGetLineOption(t *testing.T) {
	cases := []struct {
		input     string
		wantLine1 int
		wantLine2 int
		wantErr   bool
	}{
		{input: "", wantLine1: 0, wantLine2: 0, wantErr: false},
		{input: "3", wantLine1: 3, wantLine2: 0, wantErr: false},
		{input: "3-10", wantLine1: 3, wantLine2: 10, wantErr: false},
		{input: "3-", wantLine1: 0, wantLine2: 0, wantErr: true},
	}
	for _, c := range cases {
		line1, line2, err := getLineOption(c.input)
		if c.wantErr && err == nil {
			t.Errorf("'%s' wantErr %v, got %v\n", c.input, c.wantErr, err)
		}
		if line1 != c.wantLine1 || line2 != c.wantLine2 {
			t.Errorf("'%s' wantErr %v, got %v\n", c.input, c.wantErr, err)
		}
	}
}
