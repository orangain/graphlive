.PHONY: all
all: get-deps get-build-deps clean gox

.PHONY: get-deps
get-deps:
	go get golang.org/x/net/websocket

.PHONY: get-build-deps
get-build-deps:
	go get github.com/mitchellh/gox

.PHONY: clean
clean:
	rm -rf build

GOX_OPTS=-os "linux darwin windows"

.PHONY: gox
gox:
	gox $(GOX_OPTS) -output "build/{{.OS}}_{{.Arch}}/{{.Dir}}"

# For Go 1.4 or earlier
.PHONY: gox-build-toolchain
gox-build-toolchain:
	gox $(GOX_OPTS) -build-toolchain
