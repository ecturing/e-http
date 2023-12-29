build:
	go build
run:
	./ews

clean: 
	rm -f ews
all: clean build run