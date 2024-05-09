package main

import (
	"github.com/akbariandev/pacassistant/cmd/bot/commands"
	_ "go.uber.org/automaxprocs"
)

func main() {
	commands.Execute()
}
