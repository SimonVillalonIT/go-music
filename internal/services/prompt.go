package services

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/Songmu/prompter"
)

func Prompt(message string) (string, int, error) {
	search := prompter.Prompt(message, "")

	intPattern := `^\d+$`

	regex := regexp.MustCompile(intPattern)
	limit := prompter.Prompt("Enter a limit of result:", "10")

	for !regex.MatchString(limit) {
		fmt.Println("The limit must be a valid number!")
		limit = prompter.Prompt("Enter a limit of result:", "10")
	}

	limitInt, err := strconv.Atoi(limit)

	if err != nil {
		return "", 0, err
	}

	return search, limitInt, err
}
