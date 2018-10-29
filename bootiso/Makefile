NAME=bootiso
REL=latest

all: build copy

build:
	docker build --pull=true --no-cache -t $(NAME):$(REL) --build-arg rel=$(REL) .

copy:
	docker cp $$(docker create $(NAME):$(REL)):/$(NAME).iso .

.PHONY: build copy
