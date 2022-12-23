.PHONY: dev docker-build docker-run docker-tag

IMAGES=$(docker image ls | sed 0,4d | grep 'just-have-time' | awk '{print $3}')

dev:
	CompileDaemon -command="./just-have-time" -include="*.tmpl"

docker-build:
	docker image ls | grep 'just-have-time' | while read name tag id others; do if ! [ $$id = $$image_id ]; then docker image rm --force $$id; fi ; done
	docker buildx build --platform linux/amd64 -t just-have-time .

docker-run:
	docker run -p 8088:80 just-have-time

docker-lookinto:
	docker run --name just-have-time -ti -d -p 8088:80 just-have-time && docker start just-have-time && docker exec -it just-have-time /bin/sh

docker-tag:
	docker tag just-have-time gcr.io/kkchack22-just-have-time/just-have-time

docker-push:
	docker push gcr.io/kkchack22-just-have-time/just-have-time
