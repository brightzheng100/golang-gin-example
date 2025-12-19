#!/bin/bash

echo "Building server..."
env GOOS=darwin GOARCH=arm64 go build -o bin/gin-server-darwin-arm64
env GOOS=linux GOARCH=arm64 go build -o bin/gin-server-linux-arm64
env GOOS=linux GOARCH=amd64 go build -o bin/gin-server-linux-amd64
echo "Building server...DONE"

echo "Building client..."
pushd client
env GOOS=darwin GOARCH=arm64 go build -o ../bin/gin-client-darwin-arm64
env GOOS=linux GOARCH=arm64 go build -o ../bin/gin-client-linux-arm64
env GOOS=linux GOARCH=amd64 go build -o ../bin/gin-client-linux-amd64
popd
echo "Building client...DONE"
