package cmd

import (
    "errors"

	"github.com/kmaasrud/doctor"
	"github.com/kmaasrud/doctor/utils"
)

func Edit(query string) error {
	doc, err := doctor.NewDocument()
	if err != nil {
		return err
	}

    if len(doc.Chapters) > 1 {
        return errors.New("There are no chapters in this document.")
    }

	// Find the section we want to move
	matches, err := doc.FindChapterMatches(query, 0)
	if err != nil {
		return err
	}

	err = utils.OpenFile(matches[0].Path)
	if err != nil {
		return err
	}

	return nil
}
