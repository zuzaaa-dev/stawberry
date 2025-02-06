golint:
	golangci-lint run -c .golangci.yaml
.PHONY:golint

gofmt:
	gofumpt -l -w .
	goimports -w .
.PHONY:gofmt

# Frontend
npm-install:
	cd frontend && npm install

npm-run:
	cd frontend && npm run dev