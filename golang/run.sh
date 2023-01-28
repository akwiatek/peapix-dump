#!env sh

set -e

go vet main.go
go run main.go
