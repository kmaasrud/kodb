package doctor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// The character used as a separator in filenames.
const ChapterSep string = "_"

var headerRegex *regexp.Regexp = regexp.MustCompile(`(?m)^#\s+[^#\n]*`)

// Represents a chapter in the document.
type Chapter struct {
	Path  string
	Title string
	Index int
}

// Check whether this section is equal to another. Checks if their paths are equal.
func (c *Chapter) IsEqual(other *Chapter) bool {
	return c.Path == other.Path
}

// Changes the index of this section by renaming the file it represents.
func (c *Chapter) ChangeIndex(i int) error {
	c.Index = i
	newFilename := fmt.Sprintf("%02d", i) + ChapterSep + strings.Join(strings.Split(filepath.Base(c.Path), ChapterSep)[1:], "")
	newPath := filepath.Join(filepath.Dir(c.Path), newFilename)

	err := os.Rename(c.Path, newPath)
	if err != nil {
		return err
	}
	c.Path = newPath
	return nil
}

func (c *Chapter) Content() (string, error) {
    bytes, err := os.ReadFile(c.Path)
    if err != nil {
        return "", fmt.Errorf("Could not read the content of chapter %s.\n%s", c.Title, err)
    }
    return string(bytes), nil
}

// Creates a new Chapter struct from a file.
func NewChapter(path string) (Chapter, error) {
	split := strings.Split(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), ChapterSep)
    title := strings.Join(split[1:], "") // FIXME: This can fail if no underscores are in the filename

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
func PathsFromChapters(chapters []Chapter) []string {
	var paths []string
	for _, c := range chapters {
		paths = append(paths, c.Path)
	}
	return paths
}

// Finds all chapters that match the input. Returns an error if no chapters match.
// 'minus' is subtracted from the index matching statement, used in the case of looping
// over multiple inputs to match against.
func FindChapterMatches(input string, chapters []Chapter, minus int) ([]Chapter, error) {
	var matches []Chapter
	index, err := strconv.Atoi(input)
	if err != nil {
		// The input is not parsable as int, handle it as a section name
		for _, c := range chapters {
			if strings.ToLower(c.Title) == strings.ToLower(input) {
				matches = append(matches, c)
			}
		}
	} else {
		// The input is parsable as int, handle it as a section index
		// Index matching is a bit difficult, since the indices change around a lot
		// when removing multiple sections. To solve this, subtract the number of sections
		// deleted from the index matched against.
		for _, c := range chapters {
			if c.Index == index-minus {
				matches = append(matches, c)
			}
		}
	}

	if len(matches) < 1 {
		return matches, errors.New("Could not find any sections matching " + input + ".")
	}
	return matches, nil
}
