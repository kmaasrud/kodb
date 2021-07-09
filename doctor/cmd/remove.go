package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kmaasrud/doctor"
	"github.com/kmaasrud/doctor/core"
	"github.com/kmaasrud/doctor/log"
)

func Remove(inputs []string, confirm bool) error {
    msg := log.Get()
	var removeThis core.Section

	doc, err := doctor.NewDocument()
	if err != nil {
		return err
	}

    if len(doc.Chapters) < 1 {
        return errors.New("Cannot remove. This document does not contain any chapters.")
    }

	// Loop over supplied inputs and delete if they match
	for i, input := range inputs {
		matches, err := doc.FindChapterMatches(input, i)
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
            i, err := log.ChooseSection(titles, fmt.Sprintf("Found %d matches", len(matches)), "Which one do you want to delete?")
			if err != nil {
				continue
			}
            removeThis = matches[i]
		}

		// Confirmation of deletion if not already supplied on the command line
		if !confirm {
			var confirmString string
			fmt.Printf("Are you sure you want to delete \"%s\"? (y/N) ", removeThis.Title)
			fmt.Scanln(&confirmString)
			if strings.ToLower(confirmString) != "y" {
				msg.Info("Skipping deletion of \"" + removeThis.Title + "\".")
				continue
			}
		}

		// Remove the file
		err = os.Remove(removeThis.Path)
		if err != nil {
			msg.Error("Could not remove section \"" + removeThis.Title + "\". " + err.Error())
			continue
		}
		msg.Success("Deleted section \"" + removeThis.Title + "\".")

		// Decrement the sections above the removed one
		msg.Info("Reordering existing sections...")
		for j := removeThis.Index + 1; j < len(doc.Chapters); j++ {
			// Make sure we're not trying to renumber removeThis itself (if multiple sections previously shared indices)
			if doc.Chapters[j].IsEqual(removeThis) {
				continue
			}
			err = doc.Chapters[j].ChangeIndex(j - 1)
			if err != nil {
				return errors.New("Could not bump index of existing section.\n        " + err.Error())
			}
		}

		if removeThis.Index > len(doc.Chapters)-2 {
			// If the removed section has the highest index, just slice away the last element of secs
			doc.Chapters = doc.Chapters[:len(doc.Chapters)-2]
		} else {
			// Else, remove the element pertaining to this index and keep the order by reslicing
			doc.Chapters = append(doc.Chapters[:removeThis.Index], doc.Chapters[removeThis.Index+1:]...)
		}
	}

	return nil
}
