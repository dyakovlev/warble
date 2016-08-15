all: server client

server:
	go install ./...

client:
	rollup -c

clean:
	go clean
	rm -rf build/*
