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
	FieldName string
	Err       error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var s string
	if len(v) == 1 {
		return v[0].Err.Error()
	}
	for _, err := range v {
		s += "FieldName: " + err.FieldName + "\t" + "Error: " + err.Err.Error() + "\n"
	}
	return s
}

func checkValidatorTag(curValidateTag string, dt reflect.Type, ind int) *ValidationError {
	if !(dt.Field(ind).IsExported()) {
		return &ValidationError{Err: ErrValidateForUnexportedFields, FieldName: dt.Field(ind).Name}
	}
	if !utils.IsValidatorSyntaxCorrect(curValidateTag) {
		return &ValidationError{Err: ErrInvalidValidatorSyntax, FieldName: dt.Field(ind).Name}
	}
	return nil
}

func Validate(v any) error {
	errArr := make(ValidationErrors, 0)
	dt := reflect.TypeOf(v)
	if dt.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	values := reflect.ValueOf(v)

	for i := 0; i < values.NumField(); i++ {
		curValidateTag := dt.Field(i).Tag.Get(TagName)
		curFieldName := dt.Field(i).Name

		if curValidateTag == "" {
			continue
		}
		if tagErr := checkValidatorTag(curValidateTag, dt, i); tagErr != nil {
			errArr = append(errArr, *tagErr)
			continue
		}

		switch fieldValue := values.Field(i).Interface().(type) {
		case string:
			err := stringValidator.IsFieldValid(fieldValue, curValidateTag)
			if err != nil {
				errArr = append(errArr, ValidationError{Err: err, FieldName: curFieldName})
			}
		case int:
			err := intValidator.IsFieldValid(fieldValue, curValidateTag)
			if err != nil {
				errArr = append(errArr, ValidationError{Err: err, FieldName: curFieldName})
			}
		case []int:
			for _, val := range fieldValue {
				err := intValidator.IsFieldValid(val, curValidateTag)
				if err != nil {
					errArr = append(errArr, ValidationError{Err: err, FieldName: curFieldName})
				}
			}
		case []string:
			for _, val := range fieldValue {
				err := stringValidator.IsFieldValid(val, curValidateTag)
				if err != nil {
					errArr = append(errArr, ValidationError{Err: err, FieldName: curFieldName})
				}
			}
		default:
			errArr = append(errArr, ValidationError{Err: ErrUnsupportedFieldValueType, FieldName: curFieldName})
		}
	}

	if len(errArr) != 0 {
		return errArr
	}
	return nil
}
