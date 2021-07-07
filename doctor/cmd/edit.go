package cmd

import (
	"errors"

	"github.com/kmaasrud/doctor/core"
	"github.com/kmaasrud/doctor/utils"
)

func Edit(query string) error {
	rootPath, err := utils.FindDoctorRoot()
	if err != nil {
		return err
	}

	// Find all existing sections
	secs, err := utils.FindSections(rootPath)
	if err != nil {
		if _, ok := err.(*utils.NoSectionsError); !ok {
			return err
		}
		return errors.New("Could not load section list. " + err.Error())
	}

	// Find the section we want to move
	matches, err := core.FindSectionMatches(query, secs, 0)
	if err != nil {
		return err
	}

	err = utils.OpenFile(matches[0].Path)
	if err != nil {
		return err
	}

	return nil
}
