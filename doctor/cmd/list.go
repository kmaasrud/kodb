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
        msg.Info("There are no chapters in this document.")
    }

	for _, chapter := range doc.Chapters {
		fmt.Printf("%s %s\n", log.Style(fmt.Sprintf("%3d", chapter.Index), "Gray"), chapter.Title)
	}

	return nil
}
