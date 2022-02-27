.DEFAULT_GOAL := containers-push

# containers vars
fld_scripts := ./scripts
image_ver ?= dev
release_ver ?= latest
reg_name ?= 
image_options ?=

# local vars
outfile ?= _out/wordle-master

clean: containers-clean local-clean
.PHONY:clean

containers: containers-image
.PHONY:containers

local: ${outfile}
.PHONY:local

############
# containers
containers-dep-clean:
	@${fld_scripts}/docker-build-dep.sh clean
.PHONY:containers-dep-clean

containers-dep: containers-dep-clean
	@${fld_scripts}/docker-build-dep.sh
.PHONY:containers-dep

containers-image: containers-dep
	@${fld_scripts}/docker-build-image.sh ${image_ver} ${image_options}
	@$(MAKE) containers-dep-clean
.PHONY:image

containers-push: containers-image
	@${fld_scripts}/docker-aws-push-images.sh ${image_ver} ${release_ver}
.PHONY:containers-push

containers-clean: containers-dep-clean
.PHONY:containers-clean

containers-run: containers-image
	docker run --rm --name wordle -p 8080:8080 master-build:${image_ver}
.PHONY:containers-run

############
# local
${outfile}: local-build

local-fmt:
	@go fmt ./...
.PHONY:local-fmt

local-lint: local-fmt
	@golint ./...
.PHONY:local-lint

local-vet: local-fmt
	@go vet ./...
.PHONY:local-vet

local-build: local-vet
	@go build -o ${outfile} ./cmd/wordleserver
.PHONY:local-build

local-clean:
	@-go clean -cache -i -r
.PHONY:local-clean

local-dep: local-clean
	@go mod download
.PHONY:local-dep

local-deploy: local-build
	@go build -o ${outfile} ./cmd/wordleserver
.PHONY:local-deploy
