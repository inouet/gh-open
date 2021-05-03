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

(TODO)

### Download binary:

(TODO)

## Supported services

* https://github.com/ (TODO: GHE)
* https://gitlab.com/
* https://bitbucket.org/
