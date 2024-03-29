project_name = app
image_name = app:latest
postgre_image = postgres_gin_boilerplate:latest

run-local:
	go fmt ./... && gosec ./... && air app.go

docs-generate:
	swag init

requirements:
	go mod tidy

clean-packages:
	go clean -modcache

up: 
	make up-silent
	make shell

build:
	docker build -t $(image_name) .

build-no-cache:
	docker build --no-cache -t $(image_name) .

up-silent:
	make delete-container-if-exist
	make delete-postgre-if-exist
	make up-postgre
	make build
	docker run --env-file .env.dev -p 3006:3006 --name $(project_name) $(image_name) 

up-silent-prefork:
	make delete-container-if-exist
	docker run -d -p 3000:3000 --name $(project_name) $(image_name) ./app -prod

up-postgre:
	docker run --name $(postgre_image) -e POSTGRES_PASSWORD=postgrepw -e POSTGRES_DB=quizard -d -p 5500:5432 postgres

delete-postgre-if-exist:
	docker rm --force $(postgre_image)

delete-container-if-exist:
	docker stop $(project_name) || true && docker rm $(project_name) || true

shell:
	docker exec -it $(project_name) /bin/sh

stop:
	docker stop $(project_name)

start:
	docker start $(project_name)
