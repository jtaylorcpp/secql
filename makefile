generate:
	go run github.com/99designs/gqlgen generate

cli:
	mkdir -p builds
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o builds/secql cmd/*.go