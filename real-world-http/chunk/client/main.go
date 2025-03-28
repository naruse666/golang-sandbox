package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"net"
)

func main() {
	// open tcp socket
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	conn, err := dialer.Dial("tcp", "localhost:18888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// request
	request, err := http.NewRequest("GET", "http://localhost:18888", nil)
	err = request.Write(conn)
	if err != nil {
		panic(err)
	}

	// read
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, request)
	if err != nil {
		panic(err)
	}

	if resp.TransferEncoding[0] != "chunked" {
		panic("wrong transfer encoding")
	}
	for {
		sizeStr, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		// parse hex
		size, err := strconv.ParseInt(string(sizeStr[:len(sizeStr)-2]), 16, 64)
		if size == 0 {
			break
		}
		if err != nil {
			panic(err)
		}

		line := make([]byte, int(size))
		reader.Read(line)
		reader.Discard(2)
		log.Printf(" %s\n", strings.TrimSpace(string(line)))
	}
}
