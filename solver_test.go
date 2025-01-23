package main

import "testing"

func TestGcdArray(t *testing.T) {
	a := []int{4, 4, 6, 8, 16}
	solver := newSolver()
	gcd := solver.findGcdForArray(a)
	if gcd != 2 {
		t.Fatalf("Failed calculating GCD for array of int")
	}

}

func TestArrayToInt(t *testing.T) {
	a := []int{0, 1, 0}
	solver := newSolver()
	num := solver.arrayNumbersToInt(a)
	if num != 10 {
		t.Fatalf("Failed TestArrayToInt: expected: %d    actual %d", 10, num)
	}

	a = []int{1, 2, 3}
	num = solver.arrayNumbersToInt(a)
	if num != 123 {
		t.Fatalf("Failed TestArrayToInt: expected: %d    actual %d", 123, num)
	}
}
