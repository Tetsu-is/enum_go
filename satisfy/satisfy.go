package satisfy

import "errors"

var (
	ErrValueIsNotOdd = errors.New("value is not odd")
)

func ValueSatisfyOddInt(v int) error {
	if v%2 == 1 {
		return nil
	}
	return ErrValueIsNotOdd
}
