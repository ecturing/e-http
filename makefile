SERVER_BIN:=ews_server
TESTCLIENT_BIN:=client
TestOutPut:=outdata

build:
	@go build -o out/target/$(SERVER_BIN)
	@go build -o out/target/$(TESTCLIENT_BIN) ./test/Client.go

srun:
	@./out/target/$(SERVER_BIN) &
	@sleep 2

trun:
	@./out/target/$(TESTCLIENT_BIN) > $(TestOutPut)

clean: 
	@rm -f ./out/target/$(SERVER_BIN) 
	@rm -f ./out/target/$(TESTCLIENT_BIN)
	@rm -f ./out/log/*.log

testCheck: trun
	$(eval OutPut := $(shell cat $(TestOutPut)))
	@trap 'rm -f $(TestOutPut)' EXIT; \
    if echo "$(OutPut)" | grep -q "失败"; then \
        echo "failed"; \
        echo "$(OutPut)"; \
        exit 1; \
    elif echo "$(OutPut)" | grep -q "成功"; then \
        echo "success"; \
        exit 0; \
    else \
        echo "error"; \
        exit 1; \
    fi
shutdown:
	@pkill -f $(SERVER_BIN)

Stest: build srun testCheck shutdown