package main

import (
	"github.com/zer0go/ws-relay-service/cmd"
)

var Version = "development"

func main() {
	cmd.Execute(Version)
}
