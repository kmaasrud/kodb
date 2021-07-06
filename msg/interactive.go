package msg

import (
	"fmt"
	"strconv"
	"strings"
    "errors"
)

// Enters interactive mode to select among a choice of strings. Returns the chosen index
// or an error indicating that the user quit.
func ChooseSection(options []string, initMessage, choiceMessage string) (int, error) {
	var chosenIndex string
	var choice int

	Info(initMessage)
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
			Info("That is not a valid option. Please enter the number of the section you want to remove.")
		}
	}
	return choice, nil
}
