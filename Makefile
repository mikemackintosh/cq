BIN = cq

install: cmd/cq/main.go
	go install cmd/cq/main.go

clean:
	rm bin/*
