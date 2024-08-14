default: build

test:
	echo hello
	# go test ./...

build:
	go build

install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-template ~/.tflint.d/plugins
