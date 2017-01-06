#!/bin/sh

go get github.com/cesanta/docker_auth/auth_server
go get github.com/segmentio/pointer
go get github.com/zemirco/couchdb

go build
