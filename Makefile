GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

.PHONY: test cover
test:
	$(GOTEST) -v -coverprofile=coverage.out -count=1 ./...

cover:
	$(GOTEST) -v -coverprofile=coverage.out -count=1 ./...
	$(GOCOVER) -html=coverage.out
