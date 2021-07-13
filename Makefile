
default: build
NAME:=build-tools

ifeq ($(GITHUB_REF),)
GITHUB_REF := dev
endif
ifeq ($(VERSION),)
VERSION := $(shell git describe --tags --exclude "*-g*" ) built $(shell date)
endif


.PHONY: release
release: linux darwin compress


.PHONY: build
build:
	go build -o ./.bin/$(NAME) -ldflags "-X \"main.version=$(VERSION)\""  main.go


.PHONY: linux
linux:
	GOOS=linux go build -o ./.bin/$(NAME) -ldflags "-X \"main.version=$(VERSION)\""  main.go

.PHONY: darwin
darwin:
	GOOS=darwin go build -o ./.bin/$(NAME)_osx -ldflags "-X \"main.version=$(VERSION)\""  main.go

.PHONY: compress
compress:
	upx -5 -v ./.bin/build-tools

.PHONY: install
install:
	cp ./.bin/$(NAME) /usr/local/bin/

.PHONY: docker
docker:
	docker build ./ -t $(NAME)

.PHONY: test
test: docker
	command -v dgoss 2>&1 > /dev/null || test/installgoss.sh
	GOSS_FILES_PATH=test dgoss run $(NAME) sh -c "sleep 600" 
