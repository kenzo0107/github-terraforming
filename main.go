package main

import (
	"log"

	"github.com/kenzo0107/github-terraform-migration/github"
)

func main() {
	log.Println("start")
	github.Teams()
	log.Println("finished")
}
