package main

import (
    "fmt"
    "time"
)
// a simple function that returns true if a number is even
func isEven(n int) bool {
	return n%2 == 0
}

func main() {
	n := 0

  // goroutine 1
	// reads the value of n and prints true if its even
	// and false otherwise
	go func() {
		nIsEven := isEven(n)
    // we can simulate some long running step by sleeping
		// in practice, this can be some file IO operation
		// or a TCP network call
		time.Sleep(5 * time.Millisecond)
		if nIsEven {
			fmt.Println(n, " is even")
			return
		}
		fmt.Println(n, "is odd")
	}()

  // goroutine 2
	// modifies the value of n
	go func() {
		n++
	}()

	// just waiting for the goroutines to finish before exiting
	time.Sleep(time.Second)
}

