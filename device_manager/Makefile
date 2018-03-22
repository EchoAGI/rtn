all: device_manager


device_manager:
	export GOPATH=${GOPATH}:$(shell pwd) && go build -o bin/device_manager src/main.go

clean:
	go clean -i ./...
	@rm -rf bin/*
