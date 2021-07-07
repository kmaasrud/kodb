package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kmaasrud/doctor/core"
	"github.com/kmaasrud/doctor/msg"
	"github.com/kmaasrud/doctor/utils"
)

func Add(sectionName string, index int) error {
	var (
		addIndex int
		addTitle string
		addPath  string
	)

	rootPath, err := utils.FindDoctorRoot()
	if err != nil {
		return errors.New("Could not add a new section. " + err.Error())
	}

	secsDir := filepath.Join(rootPath, "secs")

	// Find all existing sections
	secs, err := utils.FindSections(rootPath)
	if _, ok := err.(*utils.NoSectionsError); ok {
		// NoSectionsError. secs directory might not exists, create it if not
		if _, existErr := os.Stat(secsDir); os.IsNotExist(existErr) {
			err := os.Mkdir(secsDir, 0755)
			if err != nil {
				return errors.New("Could not create directory 'secs'. " + err.Error())
			}
			msg.Info("Created directory 'secs'.")
		}
	} else if err != nil {
		return errors.New("Could not add a new section. " + err.Error())
	}

	if index >= 0 {
		// If index is specified, bump the index of all files above it by 1
		addIndex = index
		msg.Info("Reordering existing sections...")
		for i := index; i < len(secs); i++ {
			err := secs[i].ChangeIndex(i + 1)
			if err != nil {
				return errors.New("Could not bump index of existing section. " + err.Error())
			}
		}
	} else {
		// If no index is specified, insert the new section at the first non-occupied index
		for i, sec := range secs {
			if i < sec.Index {
				break
			}
			addIndex += 1
		}
	}

	// Title is just the supplied name, but with the first letter capitalized
	addTitle = strings.ToUpper(string(sectionName[0])) + string(sectionName[1:])

	// Paths consist of zero padded index, '_' and then the title, like this: '02_This is a section.md'
	addPath = filepath.Join(rootPath, "secs", fmt.Sprintf("%02d", addIndex)+core.SectionSep+sectionName+".md")
	err = ioutil.WriteFile(addPath, []byte("# "+addTitle+"\n\n"), 0666)
	if err != nil {
		return errors.New("Could not create new section. " + err.Error())
	}
	msg.Success(fmt.Sprintf("Created new section %s with index %d.", msg.Style(addTitle, "Bold"), addIndex))

	return nil
}
