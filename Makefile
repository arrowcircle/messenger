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
	docker build --rm=true --tag=arrowcircle/messenger -f Dockerfile .
cleanup:
	docker ps -a | grep arrowcircle/build-messenger | awk '{print $$1}' | xargs docker rm -f
	docker rmi -f arrowcircle/build-messenger
travis:
	@docker login -e $(DOCKER_EMAIL) -u $(DOCKER_USERNAME) -p $(DOCKER_PASS)
	export REPO=arrowcircle/messenger
	export TAG=`if [ "$(TRAVIS_BRANCH)" == "master" ]; then echo "latest"; else echo $(TRAVIS_BRANCH) ; fi`
	make docker
	docker tag $(REPO):$(COMMIT) $(REPO):$(TAG)
	docker tag $(REPO):$(COMMIT) $(REPO):travis-$(TRAVIS_BUILD_NUMBER)
	docker push $(REPO)

