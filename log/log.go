package log

import (
	"fmt"
	"os"
	"strings"
	"time"
    "io"
    "errors"
    "strconv"
    "sync"
)

type logger struct {
    level int
}

func (l *logger) generic(level int, symbol, text string, writer io.Writer) {
    if level >= l.level {
        split := strings.Split(text, "\n")
        for i, line := range split {
            if i == 0 {
                fmt.Fprintln(writer, fmt.Sprintf(" %s  %s", symbol, line))
            } else {
                fmt.Fprintln(writer, fmt.Sprintf("    %s", line))
            }
        }
    }
}

func (l *logger) Error(text string) {
    l.generic(3, Style("E", "Red", "Bold"), text, os.Stderr)
}

func (l *logger) Warning(text string) {
    l.generic(2, Style("W", "Yellow", "Bold"), text, os.Stderr)
}

func (l *logger) Info(text string) {
    l.generic(1, " ", Style(text, "Gray"), os.Stdout)
}

func (l *logger) Success(text string) {
    l.generic(1, Style("✓", "Green", "Bold"), text, os.Stdout)
}

func (l *logger) Debug(text string) {
    l.generic(0, "D", text, os.Stdout)
}

var log *logger
var once sync.Once
var Level int = 0

func Get() *logger {
    once.Do(func() {
        log = createLogger()
    })
    return log
}

func createLogger() *logger {
    return &logger{
        level: Level,
    }
}

// Style takes the inputted text and styles it according to
// the ANSI escape codes listed below. I should perhaps check for
// non-ANSI systems, but fuck that for now...
func Style(text string, styles ...string) string {
	code := map[string]int{
		"Red":           31,
		"Green":         32,
		"Yellow":        33,
		"Blue":          34,
		"Magenta":       35,
		"Cyan":          36,
		"Gray":          90,
		"BrightRed":     91,
		"BrightGreen":   92,
		"BrightYellow":  93,
		"BrightBlue":    94,
		"BrightMagenta": 95,
		"BrightCyan":    96,
		"Bold":          1,
		"Faint":         2,
		"Italic":        3,
		"Underline":     4,
		"Blink":         5,
		"Strike":        9,
	}

	for _, style := range styles {
		text = fmt.Sprintf("\033[%vm%v\033[0m", code[style], text)
	}

	return text
}


// Enters interactive mode to select among a choice of strings. Returns the chosen index
// or an error indicating that the user quit.
func ChooseSection(options []string, initMessage, choiceMessage string) (int, error) {
	var chosenIndex string
	var choice int

	log.Info(initMessage)
	for true {
		for i, opt := range options {
			fmt.Printf("%s %s\n", Style(fmt.Sprintf("%3d", i+1), "Gray"), opt)
		}
		fmt.Print(choiceMessage + " (q to quit) ")
		fmt.Scanln(&chosenIndex)
		index, err := strconv.Atoi(chosenIndex)
		if err == nil && index > 0 && index <= len(options) {
			choice = index-1
			break
		} else if strings.ToLower(chosenIndex) == "q" {
			return choice, errors.New("No choice done, the user quit the menu.")
		} else {
			log.Info("That is not a valid option. Please enter the number of the section you want to remove.")
		}
	}
	return choice, nil
}


func Do(doingText string, done chan struct{}) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for i := 0; ; {
		select {
		case <-ticker.C:
			i = i % 3
			dots := strings.Repeat(".", i+1) + strings.Repeat(" ", 2-i)
			fmt.Printf("\033[2K\r%s %s", Style(dots, "Gray"), doingText)
			i += 1
		case <-done:
			return
		}
	}
}

func CloseDo(done chan struct{}) {
	close(done)
	fmt.Printf("\033[2K\r")
}

// Tectonic, TeX and even Pandoc produce A LOT of noise. This function runs through each line
// of stderr and allows me to filter away lines I don't want. It also splits them into errors
// and warnings, allowing me to separate them and style them.
func CleanStderrMsg(stderr string) (string, string) {
	var warnings, errors string
	for _, line := range strings.Split(strings.TrimSuffix(stderr, "\n"), "\n") {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "! ") {
			errors += "        " + Style("TeX: ", "Bold") + strings.TrimPrefix(line, "! ") + "\n"
		} else if strings.HasPrefix(line, "[WARNING] ") {
			warnings += "        " + Style("Pandoc: ", "Bold") + strings.TrimPrefix(line, "[WARNING] ") + "\n"
		} else if strings.HasPrefix(line, "[ERROR] ") {
			errors += "        " + Style("Pandoc: ", "Bold") + strings.TrimPrefix(line, "[ERROR] ") + "\n"
		} else {
			errors += "        " + line + "\n"
		}
	}
	return warnings, errors
}
