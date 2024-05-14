BIN=baryon

baryon: clean
	go build -o ${BIN} main.go

clean:
	rm -rf ${BIN}

test:
	go test -cover ./...
