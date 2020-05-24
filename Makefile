.PHONY: clean
clean:
	rm -rf pkg/resources/pkged.go build

.PHONY: backend
backend: clean
	pkger -o pkg/resources; \
	go build -o build/watch-web; \
	cp config.yml build; \
	cp -rp scripts build

.PHONY: test
test:
	go test ./...

all: backend
