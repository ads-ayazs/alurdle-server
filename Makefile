.DEFAULT_GOAL := push

fld_scripts := ./scripts
image_ver ?= dev
release_ver ?= latest
reg_name ?= 
image_options ?=

dep-clean:
	@${fld_scripts}/docker-build-dep.sh clean
.PHONY:dep-clean

dep: dep-clean
	@${fld_scripts}/docker-build-dep.sh
.PHONY:dep

image: dep
	@${fld_scripts}/docker-build-image.sh ${image_ver} ${image_options}
	@$(MAKE) dep-clean
.PHONY:image

push: image
	@${fld_scripts}/docker-aws-push-images.sh ${image_ver} ${release_ver}
.PHONY:push

clean: dep-clean
.PHONY:clean

run-container: image
	docker run --rm --name wordle -p 8080:8080 master-build:${image_ver}
.PHONY:run-container