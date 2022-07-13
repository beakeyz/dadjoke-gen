

BIN=./bin/linux-amd64/
SERVER=server
CLIENT=client

build-server:
	go run ./build.go server 

build-client:
	go run ./build.go client
	

run-server:
	cd $(BIN) && ./$(SERVER)

run-client:
	cd $(BIN) && ./$(CLIENT)
