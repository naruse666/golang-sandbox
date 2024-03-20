package main

import (
	"fmt"
)

type company string

type Car struct {
	Name string

	providedBy company
}

func main() {
	car := &Car{
		Name:       "wagon",
		providedBy: "toyota",
	}

	fmt.Printf("Car %s providedBy %s. \n", car.Name, car.providedBy)

	lists := []map[string]interface{}{{
		"name:":   car.Name,
		"company": string(car.providedBy),
		// error
		// "company-string": String(car.providedBy),
		// "company-to-string":        ToString(&car.providedBy),
	}}

	fmt.Printf("%s", lists)
}

func ToString(p *string) (v string) {
	if p == nil {
		return v
	}

	return *p
}

func String(v string) *string {
	return &v
}
