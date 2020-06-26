#!/bin/sh

# this script is part of dxhd project, used for building and releasing dxhd

set -e

BINARY_NAME=dxhd

GIT_VERSION="$(git rev-parse --short HEAD)"
DIRTY_VERSION="$(git describe --tags --dirty --always)"
RELEASE_VERSION="$(date '+%d.%m.%Y_%H.%M')"

RELEASE_DIR="$(readlink -f "$(dirname "$0")")/releases"

GOOS=linux

build() {
	case "$1" in
		fast)
			go build -ldflags "-s -w -X main.version=$GIT_VERSION" -o "$BINARY_NAME" .
			echo "built fast build of dxhd"
			;;
		dev)
			go build -ldflags "-s -w -X main.version=$DIRTY_VERSION" -o "$BINARY_NAME" .
			echo "built dirty, developer build of dxhd"
			;;
		*)
			return 1
			;;
	esac
}

release_preconfig() {
	mkdir -p "$RELEASE_DIR/386"
	mkdir -p "$RELEASE_DIR/amd64"
	mkdir -p "$RELEASE_DIR/arm"
	mkdir -p "$RELEASE_DIR/arm64"
	return 0
}

release_build() {
	CGO_ENABLED=0

	echo "building for 386"
	GOARCH=386 go build -a -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$RELEASE_DIR/386/$BINARY_NAME"_386 .

	echo "building for amd64"
	GOARCH=amd64 go build -a -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$RELEASE_DIR/amd64/$BINARY_NAME"_amd64 .

	echo "building for arm"
	GOARCH=arm go build -a -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$RELEASE_DIR/arm/$BINARY_NAME"_arm .

	echo "building for arm64"
	GOARCH=arm64 go build -a -ldflags "-s -w -X main.version=$RELEASE_VERSION" -o "$RELEASE_DIR/arm64/$BINARY_NAME"_arm64 .

	echo

	return 0
}

release_push() {
	[ ! "$(git cherry)" = "" ] && echo "you have commits that are yet not pushed"
	git status | grep -qi '^untracked files:' && echo "you have untracked files"
	git status | grep -qi '^changes to be committed:' && echo "you have changes to be commited"

	echo 'are you sure you want to git tag and push? type "yes I want" in screaming snake case :)'
	read -r ANS
	echo

	case "$ANS" in
		YES_I_WANT)
			;;
		*)
			return 1
			;;
	esac

	LAST_TAG="$(git describe --tags --abbrev=0)"

	git tag -a "$RELEASE_VERSION" -m "$RELEASE_VERSION release"
	git push -u origin "$RELEASE_VERSION"

	return 0
}

release_commits() {
	COMMITS="$(git log --oneline "$LAST_TAG"..HEAD)"
	echo "commits since last tag ($LAST_TAG) for release page"
	echo ""
	echo "$COMMITS"
	echo ""
	echo "$COMMITS" > "$RELEASE_DIR/release_info"
}

if [ -z "$1" ]; then
	build dev
else
	case "$1" in
		fast)
			build fast
			;;
		dev)
			build dev
			;;
		release)
			release_preconfig
			release_build
			release_push
			release_commits
			;;
	esac
fi
