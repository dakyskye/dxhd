MASTERVERSION := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags --dirty --always)
DATE := $(shell date "+%d.%m.%Y_%H.%M")

fast:
	go build -ldflags "-s -w -X main.version=$(MASTERVERSION)" -o dxhd .

dev:
	go build -ldflags "-X main.version=$(VERSION)" -o dxhd .

release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-s -w -X main.version=$(DATE)" -o dxhd .
	git tag -a $(DATE) -m "release $(DATE)"
	git push origin $(DATE)
