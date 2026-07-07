package main

import (
	"log"

	"mnogo/figures/incapfigures"
)

func main() {
	println("\n инкопсуляция")
	list3, err := incapfigures.ReadFigures()
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range list3 {
		f.Print()
	}
}
