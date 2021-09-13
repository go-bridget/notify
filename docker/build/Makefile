.PHONY: all docker test

IMAGE := internal/notify:build

all: docker

docker:
	docker build --rm -t $(IMAGE) .

test:
	docker run -it --rm $(IMAGE) sh
