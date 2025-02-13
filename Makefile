# the program fails to load the WakeByAddressSingle, WakeByAddressAll and WaitOnAddress symbols from kernel32.dll. 
# View Details: https://github.com/golang/go/issues/61058
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
	set CGO_ENABLED=0
	set GOOS=windows
	set GOARCH=amd64 
	go build -race -o ${BINARY_NAME} -ldflags "-s -w" cmd/main.go

mac:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -race -o ${BINARY_NAME} -ldflags '-s -w' cmd/main.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -race -o ${BINARY_NAME} -ldflags '-s -w' cmd/main.go

run:
	./${BINARY_NAME} server:http -m release

debug:
	go run ./cmd/main.go

gorm:
	go run ./cmd/main.go gorm:gen -c model

clean:
	go clean
	rm ${BINARY_NAME}
