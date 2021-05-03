package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type buildURLFunc func(url url.URL, filePath, branch string, line string) string

// buildGithubURL build URL for Github
//   Format: https://github.com/<user>/<repos>/tree/<branch>/path/to/file.txt#L10-20
func buildGithubURL(baseURL url.URL, filePath, branch string, line string) string {
	filePath = strings.TrimLeft(filePath, "/")

	lineStr := ""
	if line != "" {
		lineStr = "L" + line
	}

	baseURL.Path = fmt.Sprintf("%s/tree/%s/%s", baseURL.Path, branch, filePath)
	baseURL.Fragment = lineStr

	return baseURL.String()
}

// buildBitbucketURL build URL for bitbucket
//   Format: https://bitbucket.org/<user>/<repos>/src/<branch>/file.txt#lines-10:20
func buildBitbucketURL(baseURL url.URL, filePath, branch string, line string) string {
	filePath = strings.TrimLeft(filePath, "/")

	lineStr := ""
	arr := strings.Split(line, "-")
	if len(arr) == 1 {
		lineStr = fmt.Sprintf("#lines-%s", arr[0])
	}
	if len(arr) == 2 {
		lineStr = fmt.Sprintf("lines-%s:%s", arr[0], arr[1])
	}
	baseURL.Path = fmt.Sprintf("%s/src/%s/%s", baseURL.Path, branch, filePath)
	baseURL.Fragment = lineStr
	return baseURL.String()
}

// buildGitlabURL build URL for gitlab
//  Format: https://gitlab.com/<user>/<repos>/-/blob/<branch>/file.txt#L10-20
func buildGitlabURL(baseURL url.URL, filePath, branch string, line string) string {
	filePath = strings.TrimLeft(filePath, "/")

	lineStr := ""
	if line != "" {
		lineStr = "L" + line
	}
	baseURL.Path = fmt.Sprintf("%s/-/blob/%s/%s", baseURL.Path, branch, filePath)
	baseURL.Fragment = lineStr

	return baseURL.String()
}

func buildURL(baseURL url.URL, path, branch, line string) (string, error) {
	buildFunc, err := getGitURLBuilder(baseURL)
	if err != nil {
		return "", err
	}

	remoteURL := buildFunc(baseURL, path, branch, line)
	return remoteURL, nil
}

func getGitURLBuilder(baseURL url.URL) (buildURLFunc, error) {
	switch baseURL.Host {
	case "bitbucket.org":
		return buildBitbucketURL, nil
	case "gitlab.com":
		return buildGitlabURL, nil
	case "github.com":
		return buildGithubURL, nil
	}
	// TODO: Support Github Enterprise
	return nil, errors.New("unknown git service")
}
