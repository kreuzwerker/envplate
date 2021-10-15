REPO := envplate
USER := kreuzwerker
FLAGS := "-X=main.build=`git rev-parse --short HEAD` -X=main.version=$(VERSION)"

.PHONY: build clean test

build:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -ldflags $(FLAGS) -o build/linux-amd64/ep bin/ep.go
	GOOS=linux GOARCH=arm go build -ldflags $(FLAGS) -o build/linux-arm/ep bin/ep.go
	GOOS=darwin GOARCH=amd64 go build -ldflags $(FLAGS) -o build/darwin-amd64/ep bin/ep.go

clean:
	rm -rf build
	
test:
	go test -cover