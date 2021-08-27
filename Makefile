VERSION ?= 0.0.12
NAME ?= "undina"
AUTHOR ?= "kevinyu0506"

build:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(AUTHOR)/$(NAME)\:$(VERSION) --push .

DEFAULT: build
