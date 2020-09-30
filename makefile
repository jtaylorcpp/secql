generate:
	go run github.com/99designs/gqlgen generate

cli:
	mkdir -p builds/cli
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o builds/cli/secql_linux server/cmd/*.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o builds/cli/secql_darwin server/cmd/*.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o builds/cli/secql.exe server/cmd/*.go

agent: generate
	mkdir -p builds/agent
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o builds/agent/secqld_linux agent/cmd/*.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o builds/agent/secqld_darwin agent/cmd/*.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o builds/agent/secqld.exe agent/cmd/*.go

build: cli agent