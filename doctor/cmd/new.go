package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kmaasrud/doctor/msg"
	"github.com/kmaasrud/doctor/utils"
)

// Default template for doctor.toml
const tomlFile string = `[meta]
title = "%s"
author = "AUTHOR"
date = "today"
`

// Creates a new document in the supplied directory. `path` is converted into an absolute path
// before being used. `defaultStructure` specifies whether to prepopulate the document with
// a default structure, namely a classic report structure.
func CreateDocumentAt(path string, defaultStructure bool) error {
	rootPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Check if specified directory exists. If not, create it. If it does, just write into it
	if fInfo, existErr := os.Stat(rootPath); os.IsNotExist(existErr) {
		err := os.Mkdir(rootPath, 0755)
		if err != nil {
			return errors.New("Could not create root directory: " + err.Error())
		}
		msg.Info("Created new directory " + rootPath + ".")
	} else if fInfo.IsDir() == false {
		message := `There is already a file with the path '%s'.
    Consider naming your document something else.`
		return errors.New(fmt.Sprintf(message, rootPath))
	}

	// Create the assets directory if it doesn't exist. If it does, just write into it
	assetPath := filepath.Join(rootPath, "assets")
	if _, err := os.Stat(assetPath); os.IsNotExist(err) {
		err := os.Mkdir(assetPath, 0755)
		if err != nil {
			return errors.New("Could not create assets directory: " + err.Error())
		}
		msg.Info("Made assets directory.")
	} else if err != nil {
		return errors.New("Could not create assets directory. " + err.Error())
	} else {
		msg.Info("The assets directory already exists, keeping it.")
	}

	// Find name of document's directory, and declare the title as the capitalized version of it
	docTitle := utils.CapitalizeFirst(filepath.Base(rootPath))
	// Create the TOML config file with rw permissions for all
	tomlPath := filepath.Join(rootPath, "doctor.toml")
	err = ioutil.WriteFile(tomlPath, []byte(fmt.Sprintf(tomlFile, docTitle)), 0666)
	if err != nil {
		return errors.New("Unable to create doctor.toml. " + err.Error())
	} else {
		msg.Info("Created config file: " + filepath.Base(tomlPath) + ".")
	}

	// Create bibliography file with rw permissions for all.
	// TODO: Consider not adding this on 'doctor new', but only when required.
	refPath := filepath.Join(assetPath, "references.bib")
	err = ioutil.WriteFile(refPath, []byte(""), 0666)
	if err != nil {
		return errors.New("Unable to create references.bib. " + err.Error())
	} else {
		msg.Info("Created bibliography file: " + filepath.Base(refPath) + ".")
	}

	// If --default flag is supplied, make default document structure
	if defaultStructure {
		// Change working dir into the newly created document
		err = os.Chdir(rootPath)
		if err != nil {
			msg.Warning("Could not navigate into your document. " + err.Error())
			return nil
		}

		msg.Info("Adding sections of default document structure...")
		for _, name := range []string{"Abstract", "Introduction", "Theory", "Method", "Results", "Discussion", "Conclusion"} {
			err = Add(name, -1)
			if err != nil {
				msg.Warning("Could not create section '" + name + "'. " + err.Error())
			}
		}
	}

	return nil
}
