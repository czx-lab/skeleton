BINARY_NAME=main
PLATFORM=Windows

ifeq ($(OS), Windows_NT)
	PLATFORM=window
	BINARY_NAME=main.exe
else
	ifeq ($(shell uname), Darwin)
		PLATFORM=mac
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
