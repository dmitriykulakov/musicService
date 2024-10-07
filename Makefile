.PHONY: run_server1 run_server2 test

DIRTEST=./test
DIRSERVER=./server
DIRSWAGGERAPI=./remoteServer

all: run_server1

run_server1: 
	sudo docker-compose -f $(DIRSWAGGERAPI)/docker-compose.yaml up --build -d
	go run $(DIRSWAGGERAPI)/main.go
	
run_server2: 
	sudo docker-compose -f $(DIRSERVER)/docker-compose.yaml up --build -d
	go run $(DIRSERVER)/main.go
	
test:
	go test $(DIRTEST) -v
	
	