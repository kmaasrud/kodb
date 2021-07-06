package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Error type returned when no sections are found from the root path.
// Mainly used for not throwing an error in 'doctor add' when there are no sections,
// but also useful to specify messages.
type NoSectionsError struct {
	errorMsg string
}

func (e *NoSectionsError) Error() string {
	return e.errorMsg
}

// The string separating the index and the name. If changed, make a due notice to users and
// either ensure backwards compatibility or have Doctor change the format automatically.
const SectionSep string = "_"

var headerRegex *regexp.Regexp = regexp.MustCompile(`(?m)^#\s+[^#\n]*`)

// Represents a section in the document.
type Section struct {
	Path  string
	Title string
	Index int
}

// Check whether this section is equal to another. Checks if their paths are equal.
func (s Section) IsEqual(other Section) bool {
	return s.Path == other.Path
}

// Changes the index of this section by renaming the file it represents.
func (s *Section) ChangeIndex(i int) error {
	s.Index = i
	newFilename := fmt.Sprintf("%02d_", i) + strings.Join(strings.Split(filepath.Base(s.Path), SectionSep)[1:], "")
	newPath := filepath.Join(filepath.Dir(s.Path), newFilename)

	err := os.Rename(s.Path, newPath)
	if err != nil {
		return err
	}
	s.Path = newPath
	return nil
}

// Creates a new Section struct from the input path.
func SectionFromPath(path string) (Section, error) {
	var title string
	split := strings.Split(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), SectionSep)
	title = strings.Join(split[1:], "")

	content, err := os.ReadFile(path)
	if err == nil {
		h1 := headerRegex.FindString(string(content))
		if h1 != "" {
			title = h1[2:]
		}
	}

	index, err := strconv.Atoi(split[0])
	if err != nil {
		return Section{}, err
	}

	return Section{path, title, index}, nil
}

// Takes a list of Section structs and returns a list of the corresponding paths.
func PathsFromSections(secs []Section) []string {
	var paths []string
	for _, sec := range secs {
		paths = append(paths, sec.Path)
	}
	return paths
}

// Finds all sections that match the input. Returns an error if no sections match.
// 'minus' is subtracted from the index matching statement, used in the case of looping
// over multiple inputs to match against.
func FindSectionMatches(input string, secs []Section, minus int) ([]Section, error) {
	var matches []Section
	index, err := strconv.Atoi(input)
	if err != nil {
		// The input is not parsable as int, handle it as a section name
		for _, sec := range secs {
			if strings.ToLower(sec.Title) == strings.ToLower(input) {
				matches = append(matches, sec)
			}
		}
	} else {
		// The input is parsable as int, handle it as a section index
		// Index matching is a bit difficult, since the indices change around a lot
		// when removing multiple sections. To solve this, subtract the number of sections
		// deleted from the index matched against.
		for _, sec := range secs {
			if sec.Index == index-minus {
				matches = append(matches, sec)
			}
		}
	}

	if len(matches) < 1 {
		return matches, errors.New("Could not find any sections matching " + input + ".")
	}
	return matches, nil
}

// Returns a slice containing core.Sections corresponding to this document
func FindSections(rootPath string) ([]Section, error) {
	var files []Section

	if _, err := os.Stat(filepath.Join(rootPath, "secs")); os.IsNotExist(err) {
		return nil, &NoSectionsError{"Empty Doctor document."}
	}

	// Walk should walk through dirs in lexical order, making sorting unecessary (luckily)
	err := filepath.Walk(filepath.Join(rootPath, "secs"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			sec, err := SectionFromPath(path)
			if err != nil {
				return err
			}
			files = append(files, sec)
		}
		return nil
	})
	if err != nil {
		return nil, err
	} else if len(files) < 1 {
		return nil, &NoSectionsError{"Empty Doctor document."}
	}

	return files, nil
}
