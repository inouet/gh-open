package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type buildURLFunc func(url url.URL, filePath, branch string, line1, line2 int) string

// buildGithubURL build URL for Github
//   Format: https://github.com/<user>/<repos>/tree/<branch>/path/to/file.txt#L10-L20
func buildGithubURL(baseURL url.URL, filePath, branch string, line1, line2 int) string {
	filePath = strings.TrimLeft(filePath, "/")

	lineStr := ""
	if line1 != 0 {
		lineStr = fmt.Sprintf("L%d", line1)
		if line2 != 0 {
			lineStr = lineStr + fmt.Sprintf("-L%d", line2)
		}
	}
	baseURL.Path = fmt.Sprintf("%s/tree/%s/%s", baseURL.Path, branch, filePath)
	baseURL.Fragment = lineStr

	return baseURL.String()
}

// buildBitbucketURL build URL for bitbucket
//   Format: https://bitbucket.org/<user>/<repos>/src/<branch>/file.txt#lines-10:20
func buildBitbucketURL(baseURL url.URL, filePath, branch string, line1, line2 int) string {
	filePath = strings.TrimLeft(filePath, "/")

	lineStr := ""
	if line1 != 0 {
		lineStr = fmt.Sprintf("lines-%d", line1)
		if line2 != 0 {
			lineStr = lineStr + fmt.Sprintf(":%d", line2)
		}
	}
	baseURL.Path = fmt.Sprintf("%s/src/%s/%s", baseURL.Path, branch, filePath)
	baseURL.Fragment = lineStr
	return baseURL.String()
}

// buildGitlabURL build URL for gitlab
//  Format: https://gitlab.com/<user>/<repos>/-/blob/<branch>/file.txt#L10-20
func buildGitlabURL(baseURL url.URL, filePath, branch string, line1, line2 int) string {
	filePath = strings.TrimLeft(filePath, "/")

	lineStr := ""
	if line1 != 0 {
		lineStr = fmt.Sprintf("L%d", line1)
		if line2 != 0 {
			lineStr = lineStr + fmt.Sprintf("-%d", line2)
		}
	}
	baseURL.Path = fmt.Sprintf("%s/-/blob/%s/%s", baseURL.Path, branch, filePath)
	baseURL.Fragment = lineStr

	return baseURL.String()
}

func buildURL(baseURL url.URL, path, branch string, line1, line2 int, urlType string) (string, error) {
	buildFunc, err := getGitURLBuilder(baseURL, urlType)
	if err != nil {
		return "", err
	}

	remoteURL := buildFunc(baseURL, path, branch, line1, line2)
	return remoteURL, nil
}

func getGitURLBuilder(baseURL url.URL, urlType string) (buildURLFunc, error) {
	host := baseURL.Host
	if urlType != "" {
		host = urlType
	}
	switch host {
	case "bitbucket.org":
		return buildBitbucketURL, nil
	case "gitlab.com":
		return buildGitlabURL, nil
	case "github.com":
		return buildGithubURL, nil
	}
	return nil, errors.New("unknown git service")
}
