BINARY_NAME=main
PLATFORM=Windows

ifeq ($(OS), Windows_NT)
	PLATFORM=Windows
	BINARY_NAME=main.go
else
	ifeq ($(shell uname), Darwin)
		PLATFORM=MacOS
	else
		PLATFORM=linux
	endif
endif

all: ${PLATFORM} run

window:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -race -o ${BINARY_NAME} -ldflags '-s -w' cmd/main.go

mac:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -race -o ${BINARY_NAME} -ldflags '-s -w' cmd/main.go

linux:	
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -race -o ${BINARY_NAME} -ldflags '-s -w' cmd/main.go

run:
	./${BINARY_NAME}

clean:
    go clean
    rm ${BINARY_NAME}
