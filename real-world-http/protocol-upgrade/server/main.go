package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func handlerUpgrade(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Connection") != "Upgrade" || r.Header.Get("Upgrade") != "MyProtocol" {
		w.WriteHeader(400)
		return
	}
	fmt.Println("Upgrade to MyProtocol")

	// get low layer socket
	c := http.NewResponseController(w)
	conn, readWriter, err := c.Hijack()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	response := http.Response{
		StatusCode: 101,
		Header:     make(http.Header),
	}
	response.Header.Set("Upgrade", "MyProtocol")
	response.Header.Set("Connection", "Upgrade")
	response.Write(conn)

	for i := 1; i <= 10; i++ {
		fmt.Fprintf(readWriter, "%d\n", i)
		fmt.Println("->", i)
		readWriter.Flush()
		recv, err := readWriter.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		fmt.Printf("<- %s", string(recv))
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/", handlerUpgrade)
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}
