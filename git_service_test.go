package main

import (
	"net/url"
	"reflect"
	"runtime"
	"testing"
)

func TestGetGitURLBuilder(t *testing.T) {
	cases := []struct {
		input string
		want  buildURLFunc
	}{
		{input: "https://github.com/", want: buildGithubURL},
		{input: "https://gitlab.com/", want: buildGitlabURL},
		{input: "https://bitbucket.org/", want: buildBitbucketURL},
		{input: "https://code.googlesource.com/", want: buildGooglesourceURL},
		{input: "https://google.com/", want: nil},
	}
	for _, c := range cases {
		u, _ := url.Parse(c.input)
		got, _ := getGitURLBuilder(*u, "")

		wantFunc := getFuncName(c.want)
		gotFunc := getFuncName(got)
		if wantFunc != gotFunc {
			t.Errorf("'%s' want %s, got %s\n", c.input, wantFunc, gotFunc)
		}
	}
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

