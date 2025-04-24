.PHONY: bench

bench: 
	go test ./... -bench=./.. -benchmem