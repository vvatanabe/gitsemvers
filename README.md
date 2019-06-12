gitsemvers
=======

## Description

Retrieve semvers from git tags

## Synopsis

```go
sv := &gitsemvers.Semvers{RepoPath: "path/to/repo"}
semvers := sv.VersionStrings()
```

## Command Line Tool

    % go get github.com/Songmu/gitsemvers/cmd/git-semvers
    % git-semvers
    v0.9.0
    ...

## Author

[Songmu](https://github.com/Songmu)
