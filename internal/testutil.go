package internal

import "testing"

// ShouldPanic returns a panic tat will report an error when
// it is not recovering from an error. Use this as in
// defer ShouldPanic(t)()
//
func ShouldPanic(t *testing.T) func() {
	return func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}
}
