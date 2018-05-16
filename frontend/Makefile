all: docker-build

clean:
	rm -f app/css/app.css app/css/app.css.map
	rm -rf .sass-cache/

build: clean
	mkdir -p app/css
	sass app.scss -I app app/css/app.css

docker-build: build
	docker build --pull=true --no-cache .

.PHONY: clean build docker-clean docker-build
