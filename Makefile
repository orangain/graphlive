.PHONY: all get-deps get-build-deps clean bindata gox package get-deploy-deps release

all: get-deps get-build-deps clean bindata gox package

get-deps:
	go get golang.org/x/net/websocket

get-build-deps:
	go get github.com/mitchellh/gox
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

clean:
	rm -rf build dist

bindata:
	go-bindata-assetfs webroot/...

GOX_OPTS=-os "linux darwin windows"
VERSION_NAME=master

gox:
	gox $(GOX_OPTS) -output "build/${VERSION_NAME}/{{.Dir}}_${VERSION_NAME}_{{.OS}}_{{.Arch}}/{{.Dir}}"

package:
	./package.sh build/${VERSION_NAME} dist/${VERSION_NAME}

get-deploy-deps:
	go get github.com/tcnksm/ghr

release:
	ghr --prerelease --replace pre-release dist/${VERSION_NAME}
