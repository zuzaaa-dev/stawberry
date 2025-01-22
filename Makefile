golint:
	golangci-lint run -c .golangci.yaml
.PHONY:golint

gofmt:
	gofumpt -l -w .
	goimports -w .
.PHONY:gofmt