.PHONY: build clean

TARGET := bin/go-vedirect

build:
	go build -o $(TARGET)

buildall:
	GOOS=darwin GOARCH=amd64 go build -o $(TARGET)-darwin-amd64
	GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o $(TARGET)-linux-mips

clean:
	$(RM) -r bin/ dist/
