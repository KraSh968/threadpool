.PHONY: all bench verify

all: test bench

bench: 
	go test ./... -bench=./.. -benchmem

test:
	go test -count=1 ./...