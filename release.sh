#!/bin/bash
set -o errexit

# define package name
WORK_DIR=/go/src/github.com/goodrain/rainbond
BASE_NAME=rainbond
IMAGE_BASE_NAME=rainbond
if [ $BUILD_IMAGE_BASE_NAME ]; 
then 
IMAGE_BASE_NAME=${BUILD_IMAGE_BASE_NAME}
fi

GO_VERSION=1.13
GATEWAY_GO_VERSION=1.13-alpine

if [ -z "$VERSION" ];then
  if [ -z "$TRAVIS_TAG" ]; then
    if [ -z "$TRAVIS_BRANCH" ]; then
      VERSION=V5.2-dev
    else
      VERSION=$TRAVIS_BRANCH-dev
    fi
  else
    VERSION=$TRAVIS_TAG
  fi
fi

buildTime=$(date +%F-%H)
git_commit=$(git log -n 1 --pretty --format=%h)

release_desc=${VERSION}-${git_commit}-${buildTime}

build::binary() {
	echo "---> build binary:$1"
	local OUTPATH=./_output/$GOOS/${BASE_NAME}-$1
	local DOCKER_PATH="./hack/contrib/docker/$1"
	HOME=`pwd`
	if [ "$1" = "eventlog" ];then
		docker build -t goodraim.me/event-build:v1 ${DOCKER_PATH}/build
		docker run --rm -e GOOS=${GOOS} -v `pwd`:${WORK_DIR} -w ${WORK_DIR} goodraim.me/event-build:v1 go build  -ldflags "-w -s -X github.com/goodrain/rainbond/cmd.version=${release_desc}"  -o ${OUTPATH} ./cmd/eventlog
	elif [ "$1" = "chaos" ];then
		docker run --rm -e GOOS=${GOOS} -v `pwd`:${WORK_DIR} -w ${WORK_DIR} -it golang:${GO_VERSION} go build -ldflags "-w -s -X github.com/goodrain/rainbond/cmd.version=${release_desc}"  -o ${OUTPATH} ./cmd/builder
	elif [ "$1" = "monitor" ];then
		docker run --rm -e GOOS=${GOOS} -v `pwd`:${WORK_DIR} -w ${WORK_DIR} -it golang:${GO_VERSION} go build -ldflags "-w -s -extldflags '-static' -X github.com/goodrain/rainbond/cmd.version=${release_desc}" -tags 'netgo static_build' -o ${OUTPATH} ./cmd/$1
	elif [ "$1" = "gateway" ];then
		docker run --rm -e GOOS=${GOOS} -v `pwd`:${WORK_DIR} -w ${WORK_DIR} -it golang:${GATEWAY_GO_VERSION} go build -ldflags "-w -s -X github.com/goodrain/rainbond/cmd.version=${release_desc}"  -o ${OUTPATH} ./cmd/$1
	else
		if [ "${ENTERPRISE}" = "true" ];then	
			echo "---> ENTERPRISE:${ENTERPRISE}"
			docker run --rm -e GOOS=${GOOS} -v `pwd`:${WORK_DIR} -w ${WORK_DIR} -it golang:${GO_VERSION} go build -ldflags "-w -s -X github.com/goodrain/rainbond/cmd.version=${release_desc} -X github.com/goodrain/rainbond/util/license.enterprise=${ENTERPRISE}"  -o ${OUTPATH} ./cmd/$1
		else
			docker run --rm -e GOOS=${GOOS} -v `pwd`:${WORK_DIR} -w ${WORK_DIR} -it golang:${GO_VERSION} go build -ldflags "-w -s -X github.com/goodrain/rainbond/cmd.version=${release_desc}"  -o ${OUTPATH} ./cmd/$1
		fi
	fi
	if [ "$GOOS" = "windows" ];then
	    mv $OUTPATH  ${OUTPATH}.exe
	fi
}

build::image() {
	local REPO_PATH="$PWD"
	pushd "./hack/contrib/docker/$1"
		echo "---> build binary:$1"
		local DOCKER_PATH="./hack/contrib/docker/$1"
		if [ "$1" = "eventlog" ];then
			docker build -t goodraim.me/event-build:v1 build
			docker run --rm -v "${REPO_PATH}":${WORK_DIR} -w ${WORK_DIR} goodraim.me/event-build:v1 go build  -ldflags "-w -s -X github.com/goodrain/rainbond/cmd.version=${release_desc}"  -o ${DOCKER_PATH}/${BASE_NAME}-$1 ./cmd/eventlog
		elif [ "$1" = "chaos" ];then
			docker run --rm -v "${REPO_PATH}":${WORK_DIR} -w ${WORK_DIR} -it golang:${GO_VERSION} go build -ldflags "-w -s -X github.com/goodrain/rainbond/cmd.version=${release_desc}"  -o ${DOCKER_PATH}/${BASE_NAME}-$1 ./cmd/builder
		elif [ "$1" = "monitor" ];then
			docker run -e CGO_ENABLED=0 --rm -v "${REPO_PATH}":${WORK_DIR} -w ${WORK_DIR} -it golang:${GO_VERSION} go build -ldflags "-w -s -X github.com/goodrain/rainbond/cmd.version=${release_desc}"  -o ${DOCKER_PATH}/${BASE_NAME}-$1 ./cmd/$1
		elif [ "$1" = "gateway" ];then
			docker run --rm -v "${REPO_PATH}":${WORK_DIR} -w ${WORK_DIR} -it golang:${GATEWAY_GO_VERSION} go build -ldflags "-w -s -X github.com/goodrain/rainbond/cmd.version=${release_desc}"  -o ${DOCKER_PATH}/${BASE_NAME}-$1 ./cmd/$1
		elif [ "$1" = "mesh-data-panel" ];then
			echo "mesh-data-panel not need build";
		else
			if [ "${ENTERPRISE}" = "true" ];then
				echo "---> ENTERPRISE:${ENTERPRISE}"
				docker run --rm -v "${REPO_PATH}":${WORK_DIR} -w ${WORK_DIR} -it golang:${GO_VERSION} go build -ldflags "-w -s -X github.com/goodrain/rainbond/cmd.version=${release_desc} -X github.com/goodrain/rainbond/util/license.enterprise=${ENTERPRISE}"  -o ${DOCKER_PATH}/${BASE_NAME}-$1 ./cmd/$1
			else
				docker run --rm -v "${REPO_PATH}":${WORK_DIR} -w ${WORK_DIR} -it golang:${GO_VERSION} go build -ldflags "-w -s -X github.com/goodrain/rainbond/cmd.version=${release_desc}"  -o ${DOCKER_PATH}/${BASE_NAME}-$1 ./cmd/$1
			fi
		fi
		echo "---> build image:$1"
		sed "s/__RELEASE_DESC__/${release_desc}/" Dockerfile > Dockerfile.release
		docker build -t "${IMAGE_BASE_NAME}/rbd-$1:${VERSION}" -f Dockerfile.release .
		if [ "$2" = "push" ];then
		    docker login -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD"
			docker push "${IMAGE_BASE_NAME}/rbd-$1:${VERSION}"
			if [ ${DOMESTIC_BASE_NAME} ];
			then
				docker tag "${IMAGE_BASE_NAME}/rbd-$1:${VERSION}" "${DOMESTIC_BASE_NAME}/${DOMESTIC_NAMESPACE}/rbd-$1:${VERSION}"
				docker login -u "$DOMESTIC_DOCKER_USERNAME" -p "$DOMESTIC_DOCKER_PASSWORD" ${DOMESTIC_BASE_NAME}
				docker push "${DOMESTIC_BASE_NAME}/${DOMESTIC_NAMESPACE}/rbd-$1:${VERSION}"
			fi
		fi	
		rm -f ./Dockerfile.release
		rm -f ./${BASE_NAME}-$1
	popd
}

build::all(){
	local build_items=(api chaos gateway monitor mq webcli worker eventlog init-probe mesh-data-panel grctl node)
	for item in "${build_items[@]}"
	do
		build::image "$item" "$1"
	done
}

case $1 in
	binary)
	    if [ "$2" = "all" ];then
			build_items=(chaos grctl node gateway monitor mq worker eventlog api init-probe)
			for item in "${build_items[@]}"
			do
				build::binary $item $1
			done
		else
		    build::binary $2	
		fi	
	;;
	*)
		if [ "$1" = "all" ];then
			build::all $2
		else
			build::image $1 $2
		fi
	;;
esac
