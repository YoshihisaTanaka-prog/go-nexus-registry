#!/bin/sh

if [ "$#" = "1" ]; then
  apt-get update && apt-get install -y gcc libsqlite3-dev
  cd "$(dirname $0)/${1}"
  go mod tidy
  CGO_ENABLED=1 go build -o server main.go
else
  echo "Please select target directory"
  exit 1
fi
