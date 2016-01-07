# Copyright (c) 2013-2016 by Michael Dvorkin. All Rights Reserved.
# Use of this source code is governed by a MIT-style license that can
# be found in the LICENSE file.

VERSION = 0.2.0
PACKAGE = github.com/michaeldv/mop/cmd/mop

run:
	go run ./cmd/mop/main.go

build:
	go build -x $(PACKAGE)

install:
	go install -x $(PACKAGE)

buildall:
	GOOS=darwin  GOARCH=amd64 go build $(GOFLAGS) -o ./bin/mop-$(VERSION)-osx-64         $(PACKAGE)
	GOOS=freebsd GOARCH=amd64 go build $(GOFLAGS) -o ./bin/mop-$(VERSION)-freebsd-64     $(PACKAGE)
	GOOS=linux   GOARCH=amd64 go build $(GOFLAGS) -o ./bin/mop-$(VERSION)-linux-64       $(PACKAGE)
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o ./bin/mop-$(VERSION)-windows-64.exe $(PACKAGE)
	GOOS=windows GOARCH=386   go build $(GOFLAGS) -o ./bin/mop-$(VERSION)-windows-32.exe $(PACKAGE)