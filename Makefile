GIT_VERSION = latest
NAMESPACE = default23
PACKAGE = otus-users
IMAGE=$(NAMESPACE)/$(PACKAGE):$(GIT_VERSION)

GOOS?=linux
GOARCH?=amd64

build:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
		go build ./cmd/$(PACKAGE)

docker-build: build
	docker build -t $(IMAGE) --build-arg PACKAGE=$(PACKAGE) .

docker-run: docker-build
	docker run -it --rm \
		-e LISTEN=$(LISTEN) \
		-e DB_URI=$(DB_URI) \
 		$(IMAGE)

docker-push: docker-build
	docker push $(IMAGE)
