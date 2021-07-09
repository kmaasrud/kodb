package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kmaasrud/doctor"
	"github.com/kmaasrud/doctor/log"
)

func Add(sectionName string, index int) error {
    msg := log.Get()
	var (
		addIndex int
		addTitle string
		addPath  string
	)

    doc, err := doctor.NewDocument()
	if err != nil {
		return err
	}

	// Find all existing sections
	if len(doc.Chapters) < 1 {
		// Source directory might not exists, create it if not
        srcDir := doc.Config.Build.SourceDir
		if _, existErr := os.Stat(srcDir); os.IsNotExist(existErr) {
			err := os.Mkdir(srcDir, 0755)
			if err != nil {
				return fmt.Errorf("Could not create directory '%s'. %s", srcDir, err.Error())
			}
			msg.Info("Created directory '" + srcDir + "'.")
		}
	} else if err != nil {
		return errors.New("Could not add a new section. " + err.Error())
	}

	if index >= 0 {
		// If index is specified, bump the index of all files above it by 1
		addIndex = index
		msg.Info("Reordering existing sections...")
		for i := index; i < len(doc.Chapters); i++ {
			err := doc.Chapters[i].ChangeIndex(i + 1)
			if err != nil {
				return errors.New("Could not bump index of existing section. " + err.Error())
			}
		}
	} else {
		// If no index is specified, insert the new section at the first non-occupied index
		for i, chapter := range doc.Chapters {
			if i < chapter.Index {
				break
			}
			addIndex += 1
		}
	}

	// Title is just the supplied name, but with the first letter capitalized
	addTitle = strings.ToUpper(string(sectionName[0])) + string(sectionName[1:])

	// Paths consist of zero padded index, '_' and then the title, like this: '02_This is a section.md'
	addPath = filepath.Join(doc.Root, doc.Config.Build.SourceDir, fmt.Sprintf("%02d", addIndex)+doctor.ChapterSep+sectionName+".md")
	err = ioutil.WriteFile(addPath, []byte("# "+addTitle+"\n\n"), 0666)
	if err != nil {
		return errors.New("Could not create new section. " + err.Error())
	}
	msg.Success(fmt.Sprintf("Created new section \"%s\" with index %d.", addTitle, addIndex))

	return nil
}
