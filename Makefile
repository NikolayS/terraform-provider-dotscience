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

test-install: install
	cp $(GOPATH)/bin/terraform-provider-dotscience example

test: test-install
	go run cmd/test/main.go

test-clean:
	rm -f example/terraform-provider-dotscience
	rm -rf example/.terraform
	rm -f example/terraform.tfstate
	rm -f example/terraform.tfstate.backup