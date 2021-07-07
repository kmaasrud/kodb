package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kmaasrud/doctor/core"
	"github.com/kmaasrud/doctor/lua"
	"github.com/kmaasrud/doctor/utils"
	"github.com/thatisuday/clapper"
)

func Stats(flags map[string]*clapper.Flag) error {
	// Wordcount
	if ok, _ := strconv.ParseBool(flags["wordcount"].Value); ok {
		var cmdArgs []string
		// Check for dependencies
		_, err := utils.CheckPath("pandoc")
		if err != nil {
			return errors.New("Wordcount failed. " + err.Error())
		}

		// Find root
		rootPath, err := utils.FindDoctorRoot()
		if err != nil {
			return errors.New("Wordcount failed. " + err.Error())
		}

		// Find source files
		secs, err := utils.FindSections(rootPath)
		if err != nil {
			return err
		}
		cmdArgs = append(cmdArgs, core.PathsFromSections(secs)...)

		wordcountFilter, err := lua.WordCountFilter()
		if err != nil {
			return err
		}
		cmdArgs = append(cmdArgs, "-L", wordcountFilter)

		// Ignore stderr
		out, err := utils.RunPandocWith(cmdArgs)
		if err != nil {
			return err
		}

		fmt.Print(out)
	}

	return nil
}
