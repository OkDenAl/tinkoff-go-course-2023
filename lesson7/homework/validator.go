package homework

import (
	"github.com/pkg/errors"
	"homework/utils"
	"homework/validators/intValidator"
	"homework/validators/stringValidator"
	"reflect"
)

const TagName = "validate"

var (
	ErrNotStruct                   = errors.New("wrong argument given, should be a struct")
	ErrInvalidValidatorSyntax      = errors.New("invalid validator syntax")
	ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")
	ErrUnsupportedFieldValueType   = errors.New("unsupported field value type")
)

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var s string
	for _, err := range v {
		s += err.Err.Error()
	}
	return s
}

func checkValidatorTag(curValidateTag string, dt reflect.Type, ind int) *ValidationError {
	if !(dt.Field(ind).IsExported()) {
		return &ValidationError{Err: ErrValidateForUnexportedFields}
	}
	if !utils.IsValidatorSyntaxCorrect(curValidateTag) {
		return &ValidationError{ErrInvalidValidatorSyntax}
	}
	return nil
}

func Validate(v any) error {
	errArr := make(ValidationErrors, 0)
	dt := reflect.TypeOf(v)
	if dt.Kind().String() != "struct" {
		return ErrNotStruct
	}
	values := reflect.ValueOf(v)

	for i := 0; i < values.NumField(); i++ {
		curValidateTag := dt.Field(i).Tag.Get(TagName)

		if curValidateTag == "" {
			continue
		}
		if tagErr := checkValidatorTag(curValidateTag, dt, i); tagErr != nil {
			errArr = append(errArr, *tagErr)
			continue
		}

		fildValueType := values.Field(i).Interface().(type)
		switch fildValueType {
		case string:
			err := stringValidator.IsFieldValid(fildValue.(string), curValidateTag)
			if err != nil {
				errArr = append(errArr, ValidationError{err})
			}
		case int:
			err := intValidator.IsFieldValid(fildValue.(int), curValidateTag)
			if err != nil {
				errArr = append(errArr, ValidationError{err})
			}
		case []int:
			for _, val := range fildValue.([]int) {
				err := intValidator.IsFieldValid(val, curValidateTag)
				if err != nil {
					errArr = append(errArr, ValidationError{err})
				}
			}
		case []string:
			for _, val := range fildValue.([]string) {
				err := stringValidator.IsFieldValid(val, curValidateTag)
				if err != nil {
					errArr = append(errArr, ValidationError{err})
				}
			}
		default:
			errArr = append(errArr, ValidationError{ErrUnsupportedFieldValueType})
		}
	}

	if len(errArr) != 0 {
		return errArr
	}
	return nil
}
