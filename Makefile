.DEFAULT_GOAL := image

fld_scripts := ./scripts
image_ver ?= dev
options ?=

dep-clean:
	${fld_scripts}/docker-build-dep.sh clean
.PHONY:dep-clean

dep: dep-clean
	${fld_scripts}/docker-build-dep.sh
.PHONY:dep

image: dep
	${fld_scripts}/docker-build-image.sh ${image_ver} ${options}
	$(MAKE) dep-clean
.PHONY:image
