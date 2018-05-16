all: docker-build

docker-clean:
	docker rmi $(NAME) &>/dev/null || true

docker-build:
	docker build --pull=true --no-cache .

.PHONY: docker-clean docker-build
