# gh-open

[![][circleci-svg]][circleci]

[circleci]: https://circleci.com/gh/inouet/gh-open/tree/master
[circleci-svg]: https://circleci.com/gh/inouet/gh-open.svg?style=shield

Open git repository website in your browser from terminal.

## Usage

Open the repo in your browser

```
$ gh-open ./inouet/gh-open

=> https://github.com/inouet/gh-open
```

Open the file in your browser

```
$ gh-open ./inouet/gh-open/main.go

=> https://github.com/inouet/gh-open/tree/(head commit hash)/main.go
```


Open the file in your browser (with line)

```
$ gh-open ./inouet/gh-open/main.go -l 10-20

=> https://github.com/inouet/gh-open/tree/(head commit hash)/main.go#L10-20
```


Open the file in your browser (with branch)

```
$ gh-open ./inouet/gh-open/main.go -b branch_name

=> https://github.com/inouet/gh-open/blob/branch_name/main.go
```

Print URL (Only print the url at the terminal)

```
$ gh-open ./inouet/gh-open -p

https://github.com/inouet/gh-open
```

## Installation

### Go user:

```
$ go get -u github.com/inouet/gh-open
```

### Mac user:


If you are on macOS and using Homebrew, you can install gh-open with the following command:

```
$ brew tap inouet/gh-open
$ brew install gh-open
```


### Download binary:

You can download the binary from [Releases Page](https://github.com/inouet/gh-open/releases).


## Configuration

Since the remote url is automatically generated from the domain output by `git config remote.origin.url`,
basically no configuration is required.

However, for example, if you are hosting GitHub Enterprise or GitLab in your own domain, it cannot judge,
so you can assist the judgment by setting as follows.

```
$ git config gh-open.urltype github.com
```

If you are using the http protocol, set as follows.

```
$ git config gh-open.protocol http
```


## Supported services

* https://github.com/
* https://gitlab.com/
* https://bitbucket.org/
