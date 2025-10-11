#!/usr/bin/sh

if [ "$#" = "1" ]; then
  cd "$(dirname $0)/${1}"
  if [ ! -e go.mod ]; then
    go mod init ${1}
  fi
  go build -o server main.go
else
  echo "Please select target directory"
  exit 1
fi
