build:
	go build -o ./jokebot .
run:
	go run .
run-background: build
	./jokebot > /dev/null 2>&1 &
