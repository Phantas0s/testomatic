package main

import (
	"log"
	"os"

	testomatic "github.com/Phantas0s/testomatic/cmd"
)

// init is called prior to main.
func init() {
	// Change the device for logging to stdout.
	log.SetOutput(os.Stdout)
}

func main() {
	testomatic.Run()
}
