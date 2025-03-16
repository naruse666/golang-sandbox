package main

import (
	"fmt"
	"net/http"
	"time"
)

func handlerChunkedResponse(w http.ResponseWriter, r *http.Request) {
	c := http.NewResponseController(w)
	for i := 1; i <= 10; i++ {
		fmt.Fprintf(w, "Chunk #%d\n", i)
		c.Flush()
		time.Sleep(500 * time.Millisecond)
	}
	c.Flush()
}

func main() {
	server := &http.Server{
		Addr: ":18888",
	}
	http.HandleFunc("/", handlerChunkedResponse)

	server.ListenAndServe()
}
