run:
	go run ./cmd/mop.go

build:
	go build ./cmd/mop.go

install:
	go install github.com/michaeldv/mop/cmd
	go get gopkg.in/gomail.v1
