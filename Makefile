docker-build: 
	docker build -t my-app .

docker-run:
	docker run --rm my-app

#.PHONY:integration-tests
#integration-tests:
#	docker-compose -f ./docker-compose-test.yml up -d
#	go test ./integration-tests/...
#	docker-compose -f ./docker-compose-test.yml down -rm
