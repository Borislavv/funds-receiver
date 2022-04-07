package main

import (
	"log"

	"gitlab.llo.su/fond/radara/cmd/radara"
)

func main() {
	// init and run app
	if err := radara.Run(); err != nil {
		log.Fatal(err)
	}
}
