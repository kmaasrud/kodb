package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kmaasrud/doctor/core"
	"github.com/kmaasrud/doctor/msg"
	"github.com/kmaasrud/doctor/utils"
)

func Remove(inputs []string, confirm bool) error {
	var removeThis core.Section

	rootPath, err := utils.FindDoctorRoot()
	if err != nil {
		return err
	}

	// Find all existing sections
	secs, err := core.FindSections(rootPath)
	if err != nil {
		if _, ok := err.(*core.NoSectionsError); !ok {
			return err
		}
		return errors.New("Could not load section list. " + err.Error())
	}

	// Loop over supplied inputs and delete if they match
	for i, input := range inputs {
		matches, err := core.FindSectionMatches(input, secs, i)
		if err != nil {
			msg.Error(err.Error())
			continue
		}

		if len(matches) == 1 {
			// Only one match, set is as the section to remove
			removeThis = matches[0]
		} else if len(matches) > 1 {
			// More than 1 match, enter interactive selection mode
            var titles []string
            for _, match := range matches {
                titles = append(titles, match.Title)
            }
            i, err := msg.ChooseSection(titles, fmt.Sprintf("Found %d matches", len(matches)), "Which one do you want to delete?")
			if err != nil {
				continue
			}
            removeThis = matches[i]
		}

		// Confirmation of deletion if not already supplied on the command line
		if !confirm {
			var confirmString string
			fmt.Printf("Are you sure you want to delete %s? (y/N) ", msg.Style(removeThis.Title, "Bold"))
			fmt.Scanln(&confirmString)
			if strings.ToLower(confirmString) != "y" {
				msg.Info("Skipping deletion of " + removeThis.Title + ".")
				continue
			}
		}

		// Remove the file
		err = os.Remove(removeThis.Path)
		if err != nil {
			msg.Error("Could not remove section " + msg.Style(removeThis.Title, "Bold") + ". " + err.Error())
			continue
		}
		msg.Success("Deleted section " + msg.Style(removeThis.Title, "Bold") + ".")

		// Decrement the sections above the removed one
		msg.Info("Reordering existing sections...")
		for j := removeThis.Index + 1; j < len(secs); j++ {
			// Make sure we're not trying to renumber removeThis itself (if multiple sections previously shared indices)
			if secs[j].IsEqual(removeThis) {
				continue
			}
			err = secs[j].ChangeIndex(j - 1)
			if err != nil {
				return errors.New("Could not bump index of existing section.\n        " + err.Error())
			}
		}

		if removeThis.Index > len(secs)-2 {
			// If the removed section has the highest index, just slice away the last element of secs
			secs = secs[:len(secs)-2]
		} else {
			// Else, remove the element pertaining to this index and keep the order by reslicing
			secs = append(secs[:removeThis.Index], secs[removeThis.Index+1:]...)
		}
	}

	return nil
}
