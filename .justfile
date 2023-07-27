exe_name := "datax"

alias s := start
alias b := build
alias t := test
alias f := fmt

default:
  @just --list

fmt:
    go fmt

build:
    go build -o ./dist/{{exe_name}} main.go
    just build-windows

start: build
    ./dist/{{exe_name}}

test:
    go test

build-windows:
    GOOS=windows GOARCH=amd64 go build -o ./dist/{{exe_name}}.exe main.go