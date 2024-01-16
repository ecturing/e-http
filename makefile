SERVER_BIN:=serevr
TESTCLIENT_BIN:=client

build:
	go build -o out/target/$(SERVER_BIN)
	go build -o out/target/$(TESTCLIENT_BIN) ./test/Client.go

srun:
	./out/target/$(SERVER_BIN)
trun:
	./out/target/$(TESTCLIENT_BIN)

clean: 
	rm -f ./out/target/$(SERVER_BIN) 
	rm -f ./out/target/$(TESTCLIENT_BIN)
	rm -f ./out/log/*.log
test:
	build
	srun
	sleep 3
	trun
	