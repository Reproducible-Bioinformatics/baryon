BIN=baryon

baryon: clean
	go build -o ${BIN} main.go

clean:
	rm -rf ${BIN}
