setup:
	go get
build: setup
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o messenger *.go
	chmod +x messenger
run:
	go run *.go
docker:
	docker build -t arrowcircle/build-messenger -f ./Dockerfile.build .
	docker run -t arrowcircle/build-messenger /bin/true
	docker cp `docker ps -q -n=1`:/messenger .
	docker build --rm=true --tag=arrowcircle/messenger:$(TRAVIS_COMMIT) -f Dockerfile .
cleanup:
	docker ps -a | grep arrowcircle/build-messenger | awk '{print $$1}' | xargs docker rm -f
	docker rmi -f arrowcircle/build-messenger
travis:
	@docker login -e $(DOCKER_EMAIL) -u $(DOCKER_USERNAME) -p $(DOCKER_PASS)
	make docker
	docker tag arrowcircle/messenger:$(TRAVIS_COMMIT) arrowcircle/messenger:$(TRAVIS_TAG)
	docker tag arrowcircle/messenger:$(TRAVIS_COMMIT) arrowcircle/messenger:latest
	docker push arrowcircle/messenger
