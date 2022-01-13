GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

.PHONY: test cover 
test:
	$(GOTEST) -v -count=1 ./...

cover:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out
