.PHONY: all
all: build

.PHONY: build
build:
	docker build -t me-box/ghrelease .

.PHONY: run
run:
	docker run -it -v $(pwd -P)/config.json:/config.json ghrelease
