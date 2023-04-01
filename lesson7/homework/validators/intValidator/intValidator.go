package intValidator

import (
	"errors"
	"homework/utils"
	"strconv"
	"strings"
)

var (
	ErrIsLessThenMin           = errors.New("int field is less then min value")
	ErrIsBiggerThenMax         = errors.New("int field is bigger then max value")
	ErrLenValidatorForIntValue = errors.New("cant apply len validator for int field")
	ErrCantFindIntInArray      = errors.New("cant find int field value in \"in\" validator tag array")
)

// IsFieldValid return nil if int field of struct is valid. And an error in opposite case.
func IsFieldValid(value int, tag string) error {
	splited := strings.Split(tag, ":")
	validator := splited[0]
	arguments := splited[1]
	switch validator {
	case "len":
		return ErrLenValidatorForIntValue
	case "min":
		minExpectedInt, _ := strconv.Atoi(arguments)
		if value < minExpectedInt {
			return ErrIsLessThenMin
		}
	case "max":
		maxExpectedInt, _ := strconv.Atoi(arguments)
		if value > maxExpectedInt {
			return ErrIsBiggerThenMax
		}
	case "in":
		if !utils.IsFieldValueInArray(strings.Split(arguments, ","), strconv.Itoa(value)) {
			return ErrCantFindIntInArray
		}
	}
	return nil
}
