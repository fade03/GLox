package interpreter

import (
	"fmt"
	"testing"
)

func Test_isEqual(t *testing.T) {
	fmt.Println(isEqual(nil, nil))
	fmt.Println(isEqual(1, 1))
	fmt.Println(isEqual(1, nil))
	fmt.Println(isEqual(nil, 1))
	fmt.Println(isEqual(1, "1"))
	fmt.Println(isEqual("1", nil))
}
