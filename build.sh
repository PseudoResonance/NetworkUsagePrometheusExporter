#!/usr/bin/env bash

GOOS=linux GOARCH=mips64 go build -o bin main.go
