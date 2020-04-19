.PHONY: clean
clean:
	rm -rf pkg/resources/pkged.go build

.PHONY: frontend
frontend: clean
	cd frontend && yarn build

.PHONY: backend
backend: frontend
	pkger -o pkg/resources; \
	go build -o build/watch-web; \
	cp config.yml build; \
	cp -rp scripts build

all: backend
