.PHONY: all get-deps get-build-deps clean gox package get-deploy-deps release

all: get-deps get-build-deps clean gox package

get-deps:
	go get golang.org/x/net/websocket

get-build-deps:
	go get github.com/mitchellh/gox

clean:
	rm -rf build dist

GOX_OPTS=-os "linux darwin windows"
VERSION_NAME=master

gox:
	gox $(GOX_OPTS) -output "build/${VERSION_NAME}/{{.OS}}_{{.Arch}}/{{.Dir}}"

package:
	./package.sh build/${VERSION_NAME} dist/${VERSION_NAME}

get-deploy-deps:
	go get github.com/tcnksm/ghr

release:
	ghr --prerelease --replace prerelease dist/${VERSION_NAME}
