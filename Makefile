
server_deps := $(wildcard *.go)
client_deps := $(wildcard client/js/*.js)

.PHONY: all clean client server

all: server client

server: $(server_deps)
	go install ./...

client: $(client_deps)
	rollup -c

clean:
	rm -rf build/*

