package validators

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsIntFieldValid(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		tag      string
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name:    "valid int min",
			value:   10,
			tag:     "min:10",
			wantErr: false,
		},
		{
			name:    "invalid int min",
			value:   0,
			tag:     "min:10",
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrIsLessThenMin)
			},
		},
		{
			name:    "valid int max",
			value:   5,
			tag:     "max:10",
			wantErr: false,
		},
		{
			name:    "invalid int max",
			value:   100,
			tag:     "max:10",
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrIsBiggerThenMax)
			},
		},
		{
			name:    "valid int in",
			value:   5,
			tag:     "in:10,5",
			wantErr: false,
		},
		{
			name:    "invalid int in",
			value:   100,
			tag:     "in:10,5",
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrCantFindIntInArray)
			},
		},
		{
			name:    "invalid int len",
			value:   100,
			tag:     "len:100",
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrLenValidatorForIntValue)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsIntFieldValid(tt.value, tt.tag)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
