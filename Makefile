all: rimi
.PHONY: clean

rimi:
	go build

clean:
	test -f rimi && rm rimi

