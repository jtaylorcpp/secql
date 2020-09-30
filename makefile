generate:
	go run github.com/99designs/gqlgen generate

packr-get:
	go get -u github.com/gobuffalo/packr/packr

packr: 
	$(shell go env GOPATH)/bin/packr

packr-clean:
	$(shell go env GOPATH)/bin/packr clean

cli: generate
	mkdir -p builds/cli
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o builds/cli/secql_linux server/cmd/*.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o builds/cli/secql_darwin server/cmd/*.go

agent: generate packr
	mkdir -p builds/agent
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o builds/agent/secqld_linux agent/cmd/*.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o builds/agent/secqld_darwin agent/cmd/*.go

build: cli agent packr-clean
