package utils

import (
	"strconv"
	"strings"
)

// IsValidatorSyntaxCorrect return true if validator syntax is correct and false in opposite case.
func IsValidatorSyntaxCorrect(tag string) bool {
	splited := strings.Split(tag, ":")
	if len(splited) != 2 {
		return false
	}
	validator := splited[0]
	arguments := splited[1]
	switch validator {
	case "len":
		fallthrough
	case "min":
		fallthrough
	case "max":
		if _, err := strconv.Atoi(arguments); err != nil {
			return false
		}
	case "in":
		if len(arguments) == 0 {
			return false
		}
	default:
		return false
	}
	return true
}

// IsFieldValueInArray checks if target value belongs to array of values.
func IsFieldValueInArray(inputArr []string, target string) bool {
	for _, val := range inputArr {
		if target == val {
			return true
		}
	}
	return false
}
