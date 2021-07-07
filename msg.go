package doctor

import (
	"fmt"
	"os"
	"strings"
	"time"
    "io"
)

func generic(level int, symbol, text string, writer io.Writer) {
    if level > logLevel {
        fmt.Fprintln(writer, fmt.Sprintf(" %s  %s", symbol, text))
    }
}

func Error(text string) {
    generic(3, style("E", "Red", "Bold"), text, os.Stderr)
}

func Warning(text string) {
    generic(2, style("W", "Yellow", "Bold"), text, os.Stderr)
}

func Info(text string) {
    generic(1, " ", style(text, "Gray"), os.Stdout)
}

func Success(text string) {
    generic(1, style("✓", "Green", "Bold"), text, os.Stdout)
}

func Debug(text string) {
    generic(0, "D", text, os.Stdout)
}

// style takes the inputted text and styles it according to
// the ANSI escape codes listed below. I should perhaps check for
// non-ANSI systems, but fuck that for now...
func style(text string, styles ...string) string {
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
