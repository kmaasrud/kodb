package core

import (
    "errors"
    "os"
    "path/filepath"

    "github.com/kmaasrud/doctor/core/conf"
    "github.com/kmaasrud/doctor/msg"
)

type Document struct {
    Root        string
    Chapters    []Chapter
    Config      *conf.Config
}

func NewDocument() (*Document, error) {
    var doc Document

    // Find root
	root, err := os.Getwd()
	if err != nil {
		msg.Error(err.Error())
	}

	for {
		if filepath.Dir(root) == root {
			return &doc, errors.New("Could not find a Doctor document.")
		} else if _, err := os.Stat(filepath.Join(root, "doctor.toml")); os.IsNotExist(err) {
			root = filepath.Dir(root)
        } else {
            break
        }
	}

    doc.Root = root

    // Find config
    doc.Config, err = conf.ConfigFromFile(filepath.Join(root, "doctor.toml"))
    if err != nil {
        return &doc, err
    }

    // Find chapters
	var chapters []Chapter
	if _, err := os.Stat(filepath.Join(root, doc.Config.Build.ChaptersDir)); os.IsNotExist(err) {
		return &doc, nil
	}

	// Walk should walk through dirs in lexical order, making sorting unecessary
	err = filepath.Walk(filepath.Join(root, doc.Config.Build.ChaptersDir), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			chapter, err := NewChapter(path)
			if err != nil {
				return err
			}
			chapters = append(chapters, chapter)
		}
		return nil
	})
	if err != nil {
		return &doc, err
	} else if len(chapters) < 1 {
		return &doc, nil
	}

    doc.Chapters = chapters

    return &doc, nil
}
