hello:
	echo "Lighthouse"

build:
	go build -o bin/lighthouse

run:
	go run main.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/lighthouse-linux-arm main.go 
	GOOS=linux GOARCH=arm64 go build -o bin/lighthouse-linux-arm64 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/lighthouse-linux-arm64 main.go

all: hello build