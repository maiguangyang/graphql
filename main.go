package main

import (
	"github.com/maiguangyang/graphql/cmd"
	"github.com/maiguangyang/graphql/events"
)

func main() {
	cmd.Execute()
}

// this is just for importing the events package and adding it to the go modules
var _ events.EventController
