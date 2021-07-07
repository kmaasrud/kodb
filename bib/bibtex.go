package bib

import (
	"github.com/nickng/bibtex"
	"os"
)

// Wrapper for the bibtex.Parse function. Takes in a path and returns a BibTex type
func ParseBibTex(path string) (*bibtex.BibTex, error) {
	f, err := os.Open("references.bib")
	if err != nil {
		return bibtex.NewBibTex(), err
	}

	bib, err := bibtex.Parse(f)
	if err != nil {
		return bib, err
	}
	
	return bib, nil
}
