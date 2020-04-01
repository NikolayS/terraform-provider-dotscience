JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)

LDFLAGS		+= -s -w
LDFLAGS		+= -X github.com/dotmesh-io/terraform-provider-dotscience/pkg/version.Version=$(VERSION)
LDFLAGS		+= -X github.com/dotmesh-io/terraform-provider-dotscience/pkg/version.Revision=$(GIT_REVISION)
LDFLAGS		+= -X github.com/dotmesh-io/terraform-provider-dotscience/pkg/version.BuildDate=$(JOBDATE)

install-release:
	@echo "++ Installing Dotscience Terraform Runner Provider"	
	go install -ldflags="$(LDFLAGS)" github.com/dotmesh-io/terraform-provider-dotscience/cmd/terraform-provider-dotscience

install:
	cd cmd/terraform-provider-dotscience && go install

install-example: install
	cp $(GOPATH)/bin/terraform-provider-dotscience example

image:
	docker build -t quay.io/dotmesh/terraform-provider-dotscience:alpha -f Dockerfile .
	#docker push quay.io/dotmesh/dotscience-tf-runner-provider:alpha

test:
	go get github.com/mfridman/tparse
	go test -json -v `go list ./... | egrep -v /tests` -cover | tparse -all -smallscreen