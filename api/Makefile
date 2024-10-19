dev:
	@command -v air > /dev/null && air || go run .

build:
	go build -o server -ldflags="-s -w" github.com/brantem/scorecard
