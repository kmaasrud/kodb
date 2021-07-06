package main

import (
	"fmt"
	"os"

	"github.com/kmaasrud/doctor/cmd"
	"github.com/kmaasrud/doctor/msg"
	"github.com/kmaasrud/doctor/utils"
	"github.com/thatisuday/clapper"
)

var VERSION = "DEV"

func main() {
	// Set up registry with commands
	registry := clapper.NewRegistry()

	rootCommand, _ := registry.Register("")
	rootCommand.AddFlag("version", "v", true, "")
	rootCommand.AddFlag("update", "u", true, "")
	rootCommand.AddFlag("help", "h", true, "")

	newCommand, _ := registry.Register("new")
	newCommand.AddArg("path", ".")
	newCommand.AddFlag("default", "d", true, "")
	newCommand.AddFlag("help", "h", true, "")

	buildCommand, _ := registry.Register("build")
	buildCommand.AddFlag("help", "h", true, "")

	addCommand, _ := registry.Register("add")
	addCommand.AddArg("name", "")
	addCommand.AddFlag("at", "i", false, "")
	addCommand.AddFlag("help", "h", true, "")

	removeCommand, _ := registry.Register("remove")
	removeCommand.AddArg("sections...", "")
	removeCommand.AddFlag("confirm", "c", true, "")
	removeCommand.AddFlag("help", "h", true, "")

	moveCommand, _ := registry.Register("move")
	moveCommand.AddArg("section", "")
	moveCommand.AddArg("to", "")
	moveCommand.AddFlag("help", "h", true, "")

	listCommand, _ := registry.Register("list")
	listCommand.AddFlag("help", "h", true, "")

	statsCommand, _ := registry.Register("stats")
	statsCommand.AddFlag("wordcount", "w", true, "")
	statsCommand.AddFlag("help", "h", true, "")

	editCommand, _ := registry.Register("edit")
	editCommand.AddArg("section", "0")
	editCommand.AddFlag("help", "h", true, "")

	// Parse commands
	command, err := registry.Parse(os.Args[1:])

	// Handle command parsing errors
	if err != nil {
		if _, ok := err.(clapper.ErrorUnknownCommand); ok {
			msg.Error("Unknown command '" + os.Args[1] + "'. Run 'doctor --help' to see a list of available commands.")
		} else if _, ok := err.(clapper.ErrorUnknownFlag); ok {
			msg.Error(fmt.Sprintf("%s. Run 'doctor --help' for further help.", utils.CapitalizeFirst(err.Error())))
		} else {
			msg.Error(err.Error())
		}
		os.Exit(1)
	}

	// Run command
	err = cmd.DoCommand(command, VERSION)

	if err != nil {
		// Some error messages are handled within the build function and just return an empty error
		if err.Error() != "" {
			msg.Error(err.Error())
		}
		os.Exit(1)
	}
}
