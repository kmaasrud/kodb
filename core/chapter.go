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
const ChapterSep string = "_"

var headerRegex *regexp.Regexp = regexp.MustCompile(`(?m)^#\s+[^#\n]*`)

// Represents a section in the document.
type Chapter struct {
	Path  string
	Title string
	Index int
}

// Check whether this section is equal to another. Checks if their paths are equal.
func (c Chapter) IsEqual(other Chapter) bool {
	return c.Path == other.Path
}

// Changes the index of this section by renaming the file it represents.
func (c *Chapter) ChangeIndex(i int) error {
	c.Index = i
	newFilename := fmt.Sprintf("%02d_", i) + strings.Join(strings.Split(filepath.Base(c.Path), ChapterSep)[1:], "")
	newPath := filepath.Join(filepath.Dir(c.Path), newFilename)

	err := os.Rename(c.Path, newPath)
	if err != nil {
		return err
	}
	c.Path = newPath
	return nil
}

// Creates a new Section struct from the input path.
func NewChapter(path string) (Chapter, error) {
	var title string
	split := strings.Split(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), ChapterSep)
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
		return Chapter{}, err
	}

	return Chapter{path, title, index}, nil
}

// Takes a list of Section structs and returns a list of the corresponding paths.
func PathsFromSections(chapters []Chapter) []string {
	var paths []string
	for _, c := range chapters {
		paths = append(paths, c.Path)
	}
	return paths
}

// Finds all sections that match the input. Returns an error if no sections match.
// 'minus' is subtracted from the index matching statement, used in the case of looping
// over multiple inputs to match against.
func FindSectionMatches(input string, secs []Chapter, minus int) ([]Chapter, error) {
	var matches []Chapter
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
func FindSections(rootPath string) ([]Chapter, error) {
	var files []Chapter

	if _, err := os.Stat(filepath.Join(rootPath, "secs")); os.IsNotExist(err) {
		return nil, &NoSectionsError{"Empty Doctor document."}
	}

	// Walk should walk through dirs in lexical order, making sorting unecessary (luckily)
	err := filepath.Walk(filepath.Join(rootPath, "secs"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			sec, err := NewChapter(path)
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
