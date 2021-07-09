package doctor

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kmaasrud/doctor/conf"
	"github.com/kmaasrud/doctor/log"
	"github.com/kmaasrud/doctor/bib"
	"github.com/kmaasrud/doctor/lua"
	"github.com/kmaasrud/doctor/utils"
)

type Document struct {
    Root        string
    Chapters    []Chapter
    Config      *conf.Config
}

func (d *Document) Build() error {
    msg := log.Get()

	pandocPath, err := utils.CheckPath("pandoc")
	if err != nil {
		return errors.New("Build failed. " + err.Error())
	}

	pdfEngine, err := utils.CheckPath(d.Config.Build.Engine)
	if err != nil {
		return errors.New("Build failed. " + err.Error())
	}

    metadataPath := filepath.Join(d.Root, "metadata.json")
    err = d.Config.WritePandocJson(metadataPath)
    if err != nil {
        return err
    }
    defer func() {
        err := os.Remove(metadataPath)
        if err != nil {
            msg.Warning("Failed to remove JSON metadata file. " + err.Error())
        }
    }()

    args := []string{
        "-s",
        "-o", filepath.Join(d.Root, d.Config.Build.Filename)+".pdf",
        "--resource-path="+strings.Join([]string{d.Root, filepath.Join(d.Root, "assets"), filepath.Join(d.Root, "secs")}, utils.ResourceSep),
        "--pdf-engine="+pdfEngine,
        "--metadata-file="+filepath.Join(d.Root, "metadata.json"),
    }

	if d.Config.Build.LuaFilters {
		for _, filter := range lua.BuildFilters() {
			args = append(args, "-L", filter)
		}
	}

	if _, err := os.Stat(filepath.Join(d.Root, "assets", d.Config.Bib.BibliographyFile)); err == nil {
		args = append(args, "-C", "--bibliography="+d.Config.Bib.BibliographyFile)

		// If a CSL style is specified, make sure it exists in assets
        // TODO: Move to own function in bib package
		if cslName := d.Config.Bib.Csl; cslName != "" {
			if val, ok := bib.Styles[cslName]; ok {
				err := os.WriteFile(filepath.Join(d.Root, "assets", cslName+".csl"), val, 0644)
				if err != nil {
					msg.Warning("Could not create CSL style, skipping it. " + err.Error())
					d.Config.Bib.Csl = ""
				}
			}
		}
	} else if os.IsNotExist(err) {
		msg.Warning("Could not find bibliography file: '" + d.Config.Bib.BibliographyFile + "'. Skipping citation processing.")
	}

    // Add all chapters
    args = append(args, PathsFromChapters(d.Chapters)...)

    cmd := exec.Command(pandocPath, args...)
    err = cmd.Run()
    if err != nil {
        return err
    }

    return nil
}

func NewDocument() (*Document, error) {
    msg := log.Get()

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
