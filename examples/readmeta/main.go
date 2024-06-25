package main

import (
	"flag"
	"log"

	"github.com/aamcrae/im"
)

func main() {
	flag.Parse()
	for _, f := range flag.Args() {
		_, err := im.ReadFromFile(f)
		if err != nil {
			log.Printf("%s: %v", f, err)
		} else {
			log.Printf("%s: successful decode", f)
		}
	}
}
