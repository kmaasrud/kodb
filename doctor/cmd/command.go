package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/kmaasrud/doctor/cmd/update"
	"github.com/kmaasrud/doctor/msg"
	"github.com/thatisuday/clapper"
)

func DoCommand(command *clapper.CommandConfig, version string) error {
	// Print help text
	Help(command.Flags["help"], command.Name)

	// Run the correct command logic
	switch command.Name {
	// Root command
	case "":
		for flag, val := range command.Flags {
			// Can discard this err, the root command's flags will always be parsable bools
			if ok, _ := strconv.ParseBool(val.Value); ok {
				switch flag {
				case "version":
					fmt.Printf("You are running Doctor %s\n", msg.Style(version, "Bold"))

				case "update":
					err := update.Update(version)
					if err != nil {
						return err
					}
				}
			}
		}

	// Create new document command
	case "new":
		var path string
		if val := command.Args["path"].Value; val != "" {
			path = val
		} else {
			path = command.Args["path"].DefaultValue
		}

		// Can discard this err, command.Flags["default"].Value will always be a parsable bool
		makeDefault, _ := strconv.ParseBool(command.Flags["default"].Value)

		err := CreateDocumentAt(path, makeDefault)
		if err != nil {
			return err
		}

	// Build the document command
	case "build":
		err := Build()
		if err != nil {
			return err
		}

	// Add a new section to the document
	case "add":
		var err error
		if command.Args["name"].Value == "" {
			return errors.New("Please supply a name for your section.")
		} else if indexString := command.Flags["at"].Value; indexString != "" {
			index, err := strconv.Atoi(indexString)
			if err != nil {
				return errors.New("Could not parse index: " + indexString + ". " + err.Error())
			}
			err = Add(command.Args["name"].Value, index)
		} else {
			err = Add(command.Args["name"].Value, -1)
			if err != nil {
				return err
			}
		}

	// Remove a section from the document
	case "remove":
		if command.Args["sections"].Value == "" {
			return errors.New("Please supply the name or index of the section(s) you want to remove.")
		}

		// Can discard this err, command.Flags["confirm"].Value will always be a parsable bool
		confirm, _ := strconv.ParseBool(command.Flags["confirm"].Value)

		err := Remove(strings.Split(command.Args["sections"].Value, ","), confirm)
		if err != nil {
			return err
		}

	// Move a section from one position to another
	case "move":
		section := command.Args["section"].Value
		toStr := command.Args["to"].Value
		if section == "" || toStr == "" {
			return errors.New("Please supply the section you want to move and the index you want to move it to.")
		}

		to, err := strconv.Atoi(toStr)
		if err != nil {
			return errors.New("Could not parse " + toStr + " as an index. Please supply a valid number.")
		}

		err = Move(section, to)
		if err != nil {
			return err
		}

	// List out all sections
	case "list":
		err := List()
		if err != nil {
			return err
		}

	case "stats":
		err := Stats(command.Flags)
		if err != nil {
			return err
		}

	case "edit":
		err := Edit(command.Args["section"].Value)
		if err != nil {
			return err
		}
	}

	return nil
}
