package conf

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Meta struct {
		Title  string      `toml:"title"`
		Author interface{} `toml:"author"`
		Date   string      `toml:"date"` // String or list of strings
	} `toml:"meta"`
	Build struct {
		Filename     string `toml:"filename" default:"document"`
		Engine       string `toml:"engine" default:"pdflatex"`
		LuaFilters   bool   `toml:"lua-filters" default:"true"`
		OutputFormat string `toml:"output-format" default:"pdf"`
	} `toml:"build"`
	Style struct {
		TwoColumn      bool     `toml:"two-column" default:"false"`
		NumberSections bool     `toml:"number-sections" default:"false"`
		DocumentClass  string   `toml:"document-class" default:"article"`
		ClassOptions   []string `toml:"class-options"` // String or list of strings
	} `toml:"style"`
	Bib struct {
		Csl                 string `toml:"csl"` // Default is Chicago MoS 17th Ed.
        BibliographyFile    string `toml:"bibliography-file" default:"references.bib"`
		BibliographyTitle   string `toml:"bibliography-title"`
		LinkCitations       bool   `toml:"link-citations" default:"true"`
		IncludeBibliography bool   `toml:"include-bibliography" default:"true"`
	} `toml:"bib"`
	Latex struct {
		Packages []string `toml:"packages"`
		Header   string   `toml:"header"`
	} `toml:"latex"`
	Html struct {
		Header string `toml:"header"`
	} `toml:"html"`
}

func ConfigFromFile(path string) (*Config, error) {
	var conf Config
	// Load the TOML bytes
	tomlBytes, err := os.ReadFile(path)
	if err != nil {
		return &conf, errors.New("Could not read config file. " + err.Error())
	}

	// Unmarshal into config struct
	err = toml.Unmarshal(tomlBytes, &conf)
	if err != nil {
		return &conf, errors.New("Could not unmarshal config file. " + err.Error())
	}

	return &conf, nil
}
