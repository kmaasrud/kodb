package cmd

import (
	"fmt"

	"github.com/kmaasrud/doctor"
	"github.com/kmaasrud/doctor/log"
)

func List() error {
    msg := log.Get()

	doc, err := doctor.NewDocument()
	if err != nil {
		return err
	}

    if len(doc.Chapters) > 1 {
        // An error message would look weird, so just give a friendly info message
        msg.Info("There are no chapters in this document.")
        return nil
    }

	for _, chapter := range doc.Chapters {
		fmt.Printf("%s %s\n", log.Style(fmt.Sprintf("%3d", chapter.Index), "Gray"), chapter.Title)
	}

	return nil
}
