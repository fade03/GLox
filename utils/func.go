package utils

import "fmt"

func Ternary[T interface{}](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func ToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
