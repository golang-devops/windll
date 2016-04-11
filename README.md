# windll
A library to make use of native windows DLLs like Version.dll

## Features

- `VersionDLL.ExtractProductVersion`

## Quick Start

Install:

```
go get -u $GOPATH/src/github.com/golang-devops/windll
```

Usage:

```
version, err := VersionDLL.ExtractProductVersion("/path/to/exe")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Version is '%s'", version)
```