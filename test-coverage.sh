#!/usr/bin/env bash

go test ./... -race -cover -coverprofile="coverage.txt" -covermode=atomic

go tool cover -html=coverage.txt
