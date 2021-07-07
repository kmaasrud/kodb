package cmd

import (
	"errors"
	"fmt"

	"github.com/kmaasrud/doctor/core"
	"github.com/kmaasrud/doctor/msg"
	"github.com/kmaasrud/doctor/utils"
)

func Move(input string, to int) error {
	var moveThis core.Section
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
	matches, err := core.FindSectionMatches(input, secs, 0)
	if err != nil {
		return err
	}

	// If multiple matches, enter interactive selection mode
	if len(matches) > 1 {
		var quit bool
		moveThis, quit = msg.ChooseSection(matches, fmt.Sprintf("Found %d matches", len(matches)), "Which one do you want to move?")
		if quit {
			return nil
		}
	} else {
		moveThis = matches[0]
	}

	// If moveThis.Index - to is positive, we have to move some sections up by one index.
	// If moveThis.Index - to is negative, we have to move some sections down by one index.
	if moveThis.Index-to > 0 {
		for i := to; i < moveThis.Index; i++ {
			err = secs[i].ChangeIndex(i + 1)
			if err != nil {
				return errors.New("Could not increase index of existing section.\n        " + err.Error())
			}
		}
	} else {
		for i := moveThis.Index + 1; i <= to; i++ {
			err = secs[i].ChangeIndex(i - 1)
			if err != nil {
				return errors.New("Could not reduce index of existing section.\n        " + err.Error())
			}
		}
	}

	prevIndex := moveThis.Index
	err = moveThis.ChangeIndex(to)
	if err != nil {
		return errors.New("Could not move section. " + err.Error())
	}

	msg.Success(fmt.Sprintf("Moved '%s' from index %d to %d", moveThis.Title, prevIndex, moveThis.Index))
	return nil
}
