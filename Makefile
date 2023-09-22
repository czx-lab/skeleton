BINARY_NAME=main

all: build run

build:
	go build -race -o ${BINARY_NAME} -ldflags '-s' cmd/main.go

run:
	go build -race -o ${BINARY_NAME} -ldflags '-s' cmd/main.go
	./${BINARY_NAME}

clean:
    go clean
    rm ${BINARY_NAME}
