VERSION := "1.0.0-RC1"
REPO := envplate
USER := kreuzwerker
TOKEN = `cat .token`
FLAGS := "-X=main.build=`git rev-parse --short HEAD` -X=main.version=$(VERSION)"

.PHONY: build clean release retract

build:
	cd bin && mkdir -p build  && gox -osarch="linux/amd64 linux/arm darwin/amd64" -ldflags $(FLAGS) -output "../build/{{.OS}}-{{.Arch}}/ep";

clean:
	rm -rf build

release:
	git tag $(VERSION) -f && git push --tags -f
	github-release release --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN)
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name ep-linux --file build/linux-amd64/ep
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name ep-linux-arm --file build/linux-arm/ep
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name ep-osx --file build/darwin-amd64/ep

retract:
	github-release delete --tag $(VERSION) -s $(TOKEN)
