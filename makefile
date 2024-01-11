SERVER_BIN:=serevr
TESTCLIENT_BIN:=client

build:
	go build -o out/exec/$(SERVER_BIN)
	go build -o out/exec/$(TESTCLIENT_BIN) ./test/Client.go

srun:
	./out/exec/$(SERVER_BIN)
trun:
	./out/exec/$(TESTCLIENT_BIN)

clean: 
	rm -f ./out/$(SERVER_BIN) 
	rm -f ./out/$(TESTCLIENT_BIN)

test:
	build
	srun
	sleep 3
	trun
	