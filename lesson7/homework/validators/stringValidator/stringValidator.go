package stringValidator

import (
	"errors"
	"homework/utils"
	"strconv"
	"strings"
)

var (
	ErrLenIsInvalid          = errors.New("len of string field is invalid")
	ErrLenIsLessThenMin      = errors.New("len of string field is less then min")
	ErrLenIsBiggerThenMax    = errors.New("len of string field is bigger then max")
	ErrCantFindStringInArray = errors.New("cant find string field value in \"in\" validator tag array")
)

// IsFieldValid return nil if string field of struct is valid. And an error in opposite case.
func IsFieldValid(value string, tag string) error {
	splited := strings.Split(tag, ":")
	validator := splited[0]
	arguments := splited[1]
	switch validator {
	case "len":
		expectedLen, _ := strconv.Atoi(arguments)
		if len(value) != expectedLen {
			return ErrLenIsInvalid
		}
	case "min":
		minExpectedLen, _ := strconv.Atoi(arguments)
		if len(value) < minExpectedLen {
			return ErrLenIsLessThenMin
		}
	case "max":
		maxExpectedLen, _ := strconv.Atoi(arguments)
		if len(value) > maxExpectedLen {
			return ErrLenIsBiggerThenMax
		}
	case "in":
		if !utils.IsFieldValueInArray(strings.Split(arguments, ","), value) {
			return ErrCantFindStringInArray
		}
	}
	return nil
}
