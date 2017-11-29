package main

import (
	"flag"
	"log"
	"os"

	"github.com/otiai10/cigger/travis"
	tern "github.com/otiai10/ternary"
)

var (
	service = flag.String("s", "travis", "Name of CI service")
	project = flag.String("p", "", "Name of project on the specified service")
	token   = flag.String("t", "", "API Token for the specified service")
	branch  = flag.String("b", "master", "Branch to build CI on specified project")
)

func init() {
	flag.Parse()
}

func main() {
	ci := construct(*service)
	if err := ci.Trigger(*project, *branch); err != nil {
		log.Fatalln(err)
	}
}

// Ci ...
type Ci interface {
	Trigger(string, ...string) error
}

func construct(srvc string) Ci {
	switch srvc {
	default:
		return travis.NewClient(tern.String(*token)(os.Getenv("TRAVIS_API_TOKEN")))
	}
}
