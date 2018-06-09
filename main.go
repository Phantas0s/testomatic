package main

import (
	"log"

	testomatic "github.com/Phantas0s/testomatic/cmd"
)

func main() {
	err := testomatic.Run()
	log.Fatal(err)
}
