package util

import (
	"fmt"
	"strings"
)

func AskConfirmation(message string) bool {
	var input string

	fmt.Println(fmt.Sprintf("%s (y/N): ", message))
	_, err := fmt.Scan(&input)
	if err != nil {
		panic(err)
	}

	input = strings.TrimSpace(input)
	input = strings.ToLower(input)

	if input == "y" || input == "yes" {
		return true
	}
	return false
}
