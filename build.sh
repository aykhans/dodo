#!/bin/bash

platforms=(
    "darwin,amd64"
    "darwin,arm64"
    "freebsd,386"
    "freebsd,amd64"
    "freebsd,arm"
    "linux,386"
    "linux,amd64"
    "linux,arm"
    "linux,arm64"
    "netbsd,386"
    "netbsd,amd64"
    "netbsd,arm"
    "openbsd,386"
    "openbsd,amd64"
    "openbsd,arm"
    "openbsd,arm64"
    "windows,386"
    "windows,amd64"
    "windows,arm64"
)

for platform in "${platforms[@]}"; do
    IFS=',' read -r build_os build_arch <<< "$platform"
    ext=""
    if [ "$build_os" == "windows" ]; then
        ext=".exe"
    fi
    GOOS="$build_os" GOARCH="$build_arch" go build -ldflags "-s -w" -o "./binaries/dodo-$build_os-$build_arch$ext"
done
