all:
	go get github.com/geetduggal/dx-download-agent
	go build -x -o dx-download-agent dx-download-agent.go

clean:
	go clean -x
	rm -vf dx-download-agent

check:
	go test -v .

install:
	go install -v .

.PHONY: all clean check install
