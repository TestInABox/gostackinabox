#Uncomment the below to enable detailed test information in the output
#OPTIONS=-test.v

test:
	go test -cover -covermode=atomic -coverprofile=.coverage ./... $(OPTIONS)

coverage: test
	go tool cover -html=.coverage
	go tool cover -func=.coverage

lint:
	golangci-lint run
