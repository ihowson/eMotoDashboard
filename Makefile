.PHONY: all deploy clean lint crossshell target


UID := $(shell id -u)
GID := $(id -g)
USERNAME := $(id -un)
GROUPNAME := $(id -gn)
CROSSHOME := "${HOME}"/crosshome

CROSS_DOCKER := docker run --rm -it -u=${UID}:${GID} -v '${PWD}':/src -v ${CROSSHOME}:/home -e HOME=/home -e BUILDER_UID=${UID} -e BUILDER_GID=${GID} -e BUILDER_USER=${USERNAME} -e BUILDER_GROUP=${GROUPNAME} -w /src crossbuild

all: dash dash-armhf

# Builds for RPi 3B+
dash-armhf: .crossbuild $(shell find . -name '*.go')
	# On MacOS: brew install arm-linux-gnueabihf-binutils  #NOT WORKING
	#GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CXX=arm-none-eabi-g++ CC=arm-none-eabi-gcc go build -o dash-armhf -tags="imguifreetype"

	# On Linux: apt install cpp-arm-linux-gnueabihf binutils-arm-linux-gnueabihf crossbuild-essential-armhf
	# Probably works better on Debian than Ubuntu
	#GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CXX=arm-linux-gnueabihf-g++ CC=arm-linux-gnueabihf-gcc go build -o dash-armhf -tags="imguifreetype" -tags=target .

	# Containerized
	mkdir -p ${CROSSHOME}
	# FIXME: hack
	# sudo chown ${USERNAME}:${GROUPNAME} ${CROSSHOME}
	# FIXME: hack
	chmod 777 ${CROSSHOME}
	${CROSS_DOCKER} bash -c 'PKG_CONFIG=arm-linux-gnueabihf-pkg-config GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CXX=arm-linux-gnueabihf-g++ CC=arm-linux-gnueabihf-gcc LD=arm-linux-gnueabihf-ld go build -o dash-armhf -tags="imguifreetype" -tags=target .'

crossshell: .crossbuild
	mkdir -p ${CROSSHOME}
	# FIXME: hack
	# sudo chown ${USERNAME}:${GROUPNAME} "${HOME}"/crosshome
	# FIXME: hack
	chmod 777 ${CROSSHOME}
	${CROSS_DOCKER} /bin/bash

dash: $(shell find . -name '*.go')
	go build -o dash -tags="imguifreetype" -tags=mock .

deploy: dash-armhf
	rsync -avxP dash-armhf emoto:dash
	rsync -avxP kiosk.sh emoto:kiosk.sh
	rsync -avxP assets/ emoto:assets/
	ssh emoto sudo service lightdm restart

.crossbuild: Dockerfile
	docker build -t crossbuild .
	touch .crossbuild

run: dash
	DISPLAY=:0.0 ./dash
	# ./dash

clean:
	go clean
	rm -f dash dash-armhf
	docker rm -f crossbuild || true
	rm -rf ${CROSSHOME}

lint:
	# Linux (bash)
	# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.46.2

	golangci-lint run
