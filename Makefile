
build-go-sugar:
	go build -o ./bin/go-sugar ./cmd/go-sugar

build-sugar-check:
	go build -o ./bin/sugar-check ./cmd/sugar-check

build-sugar-require:
	go build -o ./bin/sugar-require ./cmd/sugar-require

build-sugars: build-sugar-check build-sugar-require

build: build-sugars build-go-sugar
