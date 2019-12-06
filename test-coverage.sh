#!/usr/bin/env bash

go test ./... -cover -coverprofile="coverage.txt" -covermode=atomic

go tool cover -html=coverage.txt
