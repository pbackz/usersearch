#!/usr/bin/env bash
npm install
mkdir gosearch && cp -a go.mod go.sum gosearch.go data.json gosearch/ && go build && PATH=$PATH:$(pwd)