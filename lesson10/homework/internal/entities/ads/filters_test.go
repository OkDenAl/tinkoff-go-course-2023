package ads

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type FiltersIsValidTestCase struct {
	Name   string
	In     any
	Expect bool
}

type IsValidDateTest struct {
	Date string
}
type IsValidStatusTest struct {
	Stat Status
}
type IsValidAuthorIdTest struct {
	Id string
}

func TestFiltersIsValid(t *testing.T) {
	tests := []FiltersIsValidTestCase{
		{
			Name:   "correct date",
			In:     IsValidDateTest{Date: "2023-04-28"},
			Expect: true,
		},
		{
			Name:   "empty date",
			In:     IsValidDateTest{Date: ""},
			Expect: true,
		},
		{
			Name:   "incorrect month",
			In:     IsValidDateTest{Date: "2023-20-28"},
			Expect: false,
		},
		{
			Name:   "incorrect day",
			In:     IsValidDateTest{Date: "2024-02-12"},
			Expect: false,
		},
		{
			Name:   "incorrect year",
			In:     IsValidDateTest{Date: "2023-02-40"},
			Expect: false,
		},
		{
			Name:   "incorrect date",
			In:     IsValidDateTest{Date: "date"},
			Expect: false,
		},
		{
			Name:   "incorrect date",
			In:     IsValidDateTest{Date: "year-month-day"},
			Expect: false,
		},
		{
			Name:   "incorrect date len",
			In:     IsValidDateTest{Date: "2023-04-28-1"},
			Expect: false,
		},
		{
			Name:   "correct status published",
			In:     IsValidStatusTest{Stat: Published},
			Expect: true,
		},
		{
			Name:   "correct status unpublished",
			In:     IsValidStatusTest{Stat: Unpublished},
			Expect: true,
		},
		{
			Name:   "incorrect status",
			In:     IsValidStatusTest{Stat: Status("incorrect")},
			Expect: false,
		},
		{
			Name:   "correct author id",
			In:     IsValidAuthorIdTest{Id: "10"},
			Expect: true,
		},
		{
			Name:   "correct empty author id",
			In:     IsValidAuthorIdTest{Id: ""},
			Expect: true,
		},
		{
			Name:   "incorrect author id",
			In:     IsValidAuthorIdTest{Id: "hahaha"},
			Expect: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			var got bool
			switch tc.In.(type) {
			case IsValidDateTest:
				got = isDateValid(tc.In.(IsValidDateTest).Date)
			case IsValidStatusTest:
				got = isStatusValid(tc.In.(IsValidStatusTest).Stat)
			case IsValidAuthorIdTest:
				got = isAuthorIdValid(tc.In.(IsValidAuthorIdTest).Id)
			}
			assert.Equal(t, got, tc.Expect)
		})
	}
}

func BenchmarkIsAuthorIdValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = isAuthorIdValid(strconv.Itoa(i))
	}
}

func BenchmarkIsAuthorIdValid2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = isAuthorIdValid2(strconv.Itoa(i))
	}
}

func BenchmarkIsAuthorIdValid3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = isAuthorIdValid3(strconv.Itoa(i))
	}
}
