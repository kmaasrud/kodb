package conf

import (
	"encoding/json"
	"errors"
	"os"
)

type PandocConfig struct {
	Title               string      `json:"title,omitempty"`
	Author              interface{} `json:"author,omitempty"`
	Date                string      `json:"date"`
	DocumentClass       string      `json:"documentclass"`
	ClassOptions        []string    `json:"classoption,omitempty"`
	NumberSections      bool        `json:"numbersections"`
	BiblographyTitle    string      `json:"reference-section-title,omitempty"`
	Csl                 string      `json:"csl,omitempty"`
	LinkCitations       bool        `json:"link-citations"`
	HeaderIncludes      string      `json:"header-includes,omitempty"`
	SupressBibliography bool        `json:"supress-bibliography"`
}

func (c *Config) WritePandocJson(path string) error {
	pandocConf := PandocConfig{
		Title:               c.Meta.Title,
		Author:              c.Meta.Author,
		DocumentClass:       c.Style.DocumentClass,
		ClassOptions:        c.Style.ClassOptions,
		NumberSections:      c.Style.NumberSections,
		BiblographyTitle:    c.Bib.BibliographyTitle,
		Csl:                 c.Bib.Csl,
		LinkCitations:       c.Bib.LinkCitations,
		SupressBibliography: !c.Bib.IncludeBibliography,
	}

	// Handle date
	switch c.Meta.Date {
	case "today", "now", "present":
		pandocConf.Date = "\\today"
	default:
		pandocConf.Date = c.Meta.Date
	}

	// Handle output specific options
	switch c.Build.OutputFormat {
	case "html":
		pandocConf.HeaderIncludes = c.Html.Header
	case "pdf":
		if len(c.Latex.Packages) > 0 {
			for _, pack := range c.Latex.Packages {
				if string(pack[0]) == "[" || string(pack[0]) == "{" {
					// Some extra options are specified
					pandocConf.HeaderIncludes += "\\usepackage" + pack + "\n"
				} else {
					// Just a normal package name
					pandocConf.HeaderIncludes += "\\usepackage{" + pack + "}\n"
				}
			}
		}
		pandocConf.HeaderIncludes += c.Latex.Header
	}

	// Handle two-column
	if c.Style.TwoColumn {
		pandocConf.ClassOptions = append(pandocConf.ClassOptions, "twocolumn")
	}

	// Marshal config struct into JSON
	jsonBytes, err := json.Marshal(&pandocConf)
	if err != nil {
		return errors.New("Could not marshal metadata into JSON. " + err.Error())
	}

	// Write the JSON to path
	err = os.WriteFile(path, jsonBytes, 0644)
	if err != nil {
		return errors.New("Could not write JSON file. " + err.Error())
	}

	return nil
}
