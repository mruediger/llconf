package main

import (
	"log"
)

var bootstrap = &Command {
	Name: "bootstrap",
	Usage: "bootstrap [url]",
	Run: runBootstrap,
}


func runBootstrap(args []string, logi, loge *log.Logger) {

}
