.PHONY: bench verify

bench: 
	go test ./... -bench=./.. -benchmem

verify:
	go test ./... -bench=./.. -benchmem -trackGoroCount