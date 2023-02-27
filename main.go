package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	x := []int{1, 2, 3, 4}
	y := make([]int, 4)
	num := copy(y, x) // copy x -> y
	fmt.Println(y, num)

	m := map[string]int{
		"hello": 3,
		"world": 0,
		"!!":    2,
	}

	v, ok := m["hello"]
	fmt.Println(v, ok)
	// comma ok idiom
	if _, ok := m["yes"]; !ok {
		fmt.Println("key yes does not exist.")
	}
	result, remainder, err := div(128, 3)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result, remainder)
}

func div(numerator, denominator int) (int, int, error) {
	if denominator == 0 {
		return 0, 0, errors.New("can't div 0.")
	}
	return numerator / denominator, numerator % denominator, nil
}
