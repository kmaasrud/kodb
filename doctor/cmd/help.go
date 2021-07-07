package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/thatisuday/clapper"
)

func Help(helpFlag *clapper.Flag, cmdName string) {
	if ok, _ := strconv.ParseBool(helpFlag.Value); ok {
		fmt.Println("\n" + helpText[cmdName] + "\n")
		os.Exit(0)
	}
}

var helpText = map[string]string{
	"": `Doctor is a command line tool and environment for building scientific documents.

Usage:
    doctor <command> [arguments]

Commands:
    add             Add a section
    build           Build the document
    edit            Open a section in an editor
    list            List the document structure
    move            Move a section
    new             Creates a new document
    remove          Remove a section
    stats           Show document statistics

Flags:
    --help          Shows this text
    --update        Checks if an update is available and offers to download it
    --version       Prints the current version

For more help on each sub-command, run 'doctor <command> --help'.`,

	// ---
	"add": `Add a new section to your Doctor document.

Usage:
    doctor add <section name> [options]

Options:
    --at            Specify which index the new section should have`,

	// ---
	"build": `Builds the current document according to the instructions
in doctor.toml.

Usage:
    doctor build`,

	// ---
	"list": `List all sections in the document, along with their index.

Usage:
    doctor list`,

	// ---
	"move": `Move a section to a new location in the document.

Usage:
    doctor move <section> <destination>

<section> can be either the name or index of an existing section. <destination>
is the index you want to move the selected section to.

For an overview of the current document structure, you can run 'doctor list'.`,

	// ---
	"new": `Create a new Doctor document in the specified location.

Usage:
	doctor new [<directory>] [options]

[<directory>] is an optional argument specifying either an existing directory,
or the name of a new one. If the argument is omitted, Doctor will create a
document in the current directory.

Options:
    --default       Initialize the document with a standard report type
                    document structure. Very limited at the moment.`,

	// ---
	"remove": `Removes a section, or multiple sections, from the document.

Usage:
    doctor remove <section(s)...> [options]

The section(s) are specified either by their name or their index. For an 
overview, you can run 'doctor list' to list all sections and their indices.

Options:
    --confirm       Do not ask for confirmation before deleting.
                    Use this option with caution...`,

	// ---
	"stats": `Show statistics about the current Doctor document.

Usage:
    doctor stats [options]

Options:
    --wordcount     Show the number of words and characters in the document.`,

	// ---
	"edit": `Open a section in your preferred editor. 

Usage:
    doctor edit <section>

Uses the environment variable EDITOR if it is available. If not, 'xdg-open',
'open' or 'start' will be used, depending on the operating system.`,
}
