GOMOD=on

default: all

all: build start

build:
	echo -n "build server..."	
	CGO_ENABLED=1 go build -o ./bin/server ./cmd/*.go
start:
	./bin/server