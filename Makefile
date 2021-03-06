BIN = cq

install: cmd/cq/main.go
	go install cmd/cq/main.go

build:
	mkdir -p bin/
	go build -o bin/$(BIN) cmd/cq/main.go

clean:
	rm bin/*
