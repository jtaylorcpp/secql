generate:
	go run github.com/99designs/gqlgen generate

cli:
	mkdir -p builds
	go build -o builds/secql cmd/*.go