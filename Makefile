DOCKER_IMAGE = abema/github-actions-merger

.PHONY: docker-build
docker-build: TAG := latest
docker-build:
	docker build -t ${DOCKER_IMAGE}:${TAG} .

.PHONY: docker-push
docker-push: TAG := latest
docker-push: docker-build
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
