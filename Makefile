TOKEN = `cat .token`
REPO := envplate
USER := kreuzwerker
VERSION := "v0.0.4"

build: darwin linux

darwin:
	mkdir -p out/darwin
	GOOS=darwin go build -o out/darwin/ep -ldflags "-X main.build `git rev-parse --short HEAD`" bin/envplate.go
linux:
	mkdir -p out/linux
	GOOS=linux go build -o out/linux/ep -ldflags "-X main.build `git rev-parse --short HEAD`" bin/envplate.go
clean:
	rm -rf out

release:
	github-release release --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN)
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name ep-osx --file out/darwin/ep
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name ep-linux --file out/linux/ep

test:
	DATABASE=db.example.com MODE=debug go test -cover
