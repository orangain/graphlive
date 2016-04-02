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

.PHONY: gox
gox:
	gox -os "linux darwin windows" -output "build/{{.OS}}_{{.Arch}}/{{.Dir}}"
