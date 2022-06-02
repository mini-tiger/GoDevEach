package test

import (
	"test_ce/business"
	"testing"
)

func TestSum(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	expected := 15
	actual := business.Sum(numbers)
	//actual1 := business.Sum1(numbers)
	if actual != expected {
		t.Errorf("Expected the sum of %v to be %d but instead got %d!", numbers, expected, actual)
	}

}
func TestSum1(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	expected := 15
	//actual := business.Sum(numbers)
	actual1 := business.Sum1(numbers)

	if actual1 != expected {
		t.Errorf("Expected the sum1 of %v to be %d ", numbers, expected)
	}
}
