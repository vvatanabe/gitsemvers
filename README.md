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

    % go get github.com/vvatanabe/gitsemvers/cmd/git-semvers
    % git-semvers
    v0.0.1
    ...

## Origin Author

[Songmu](https://github.com/Songmu)

## Author

[vvatanabe](https://github.com/vvatanabe)
