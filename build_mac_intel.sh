#!/bin/bash

# Get current directory
cd "$(dirname "$0")"

prjPath="$(pwd)/"
binPath="${prjPath}bin/"

echo "Project Dir:${prjPath}"
echo "Target Dir:${binPath}"

# export GOARCH=arm64
# export GOOS=darwin3
export GOPATH="${prjPath}../../"

go build -o "${binPath}gxe-intel"
cp -f "${prjPath}conf.yaml" "${binPath}"
cp -rf "${prjPath}template" "${binPath}"

echo "build is completed."