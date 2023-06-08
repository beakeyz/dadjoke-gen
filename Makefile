
DEBUG_ARGS := -hostname=localhost

.PHONY: build
build:
	go run ./ build

.PHONY: debug
debug:
	./bin/linux-amd64/dadjoke-gen $(DEBUG_ARGS)
