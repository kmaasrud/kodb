package bib

import _ "embed"

//go:embed styles/apa.csl
var apa []byte

//go:embed styles/ieee.csl
var ieee []byte

//go:embed styles/harvard.csl
var harvard []byte

//go:embed styles/nature.csl
var nature []byte

//go:embed styles/vancouver.csl
var vancouver []byte

var Styles = map[string][]byte{
	"apa":       apa,
	"ieee":      ieee,
	"harvard":   harvard,
	"nature":    nature,
	"vancouver": vancouver,
}
