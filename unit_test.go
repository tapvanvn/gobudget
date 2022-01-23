package gobudget_test

import (
	"fmt"
	"math"
	"testing"
)

func TestDivision(t *testing.T) {
	timestamp1 := int64(1640354986 - 36000)
	test1 := int64(timestamp1) / 86400
	testFloat := float64(timestamp1) / 86400

	fmt.Println(timestamp1, testFloat)

	if test1 != 18986 || int64(math.Floor(testFloat)) != test1 {
		t.Fail()
	}
}
