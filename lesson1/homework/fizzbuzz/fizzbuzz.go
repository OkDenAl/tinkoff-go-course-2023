package fizzbuzz

import "strconv"

func FizzBuzz(i int) string {
	var ans string
	switch {
	case i%15 == 0:
		ans = "FizzBuzz"
	case i%3 == 0:
		ans = "Fizz"
	case i%5 == 0:
		ans = "Buzz"
	default:
		ans = strconv.Itoa(i)
	}
	return ans
}
