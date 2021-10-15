VERSION := "0.1.0"
REPO := envplate
USER := kreuzwerker
TOKEN = `cat .token`
FLAGS := "-X=main.build=`git rev-parse --short HEAD` -X=main.version=$(VERSION)"

.PHONY: build clean test release retract

build:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -ldflags $(FLAGS) -o build/linux-amd64/ep bin/ep.go
	GOOS=linux GOARCH=arm go build -ldflags $(FLAGS) -o build/linux-arm/ep bin/ep.go
	GOOS=darwin GOARCH=amd64 go build -ldflags $(FLAGS) -o build/darwin-amd64/ep bin/ep.go

clean:
	rm -rf build

release:
	git tag $(VERSION) -f && git push --tags -f
	github-release release --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN)
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name ep-osx --file build/darwin/ep
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name ep-linux --file build/linux-amd64/ep
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name ep-linux-arm --file build/linux-arm/ep
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name ep-osx --file build/darwin-amd64/ep

retract:
	github-release delete --tag $(VERSION) -s $(TOKEN)
