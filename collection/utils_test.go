package collection

import (
	"testing"
)

func TestFilter(t *testing.T) {
	// Test filtering even numbers
	numbers := []int{1, 2, 3, 4, 5, 6}
	evens := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})

	if len(evens) != 3 {
		t.Errorf("Expected 3 even numbers, got %d", len(evens))
	}

	for _, n := range evens {
		if n%2 != 0 {
			t.Errorf("Expected even number, got %d", n)
		}
	}

	// Test filtering strings by length
	strings := []string{"a", "ab", "abc", "abcd"}
	longStrings := Filter(strings, func(s string) bool {
		return len(s) > 2
	})

	if len(longStrings) != 2 {
		t.Errorf("Expected 2 strings with length > 2, got %d", len(longStrings))
	}

	for _, s := range longStrings {
		if len(s) <= 2 {
			t.Errorf("Expected string with length > 2, got %s", s)
		}
	}
}

func TestFirst(t *testing.T) {
	// Test finding first even number
	numbers := []int{1, 3, 5, 6, 7, 8}
	firstEven := First(numbers, func(n int) bool {
		return n%2 == 0
	})

	if firstEven != 6 {
		t.Errorf("Expected first even number to be 6, got %d", firstEven)
	}

	// Test finding first string with specific prefix
	strings := []string{"foo", "bar", "baz", "qux"}
	firstWithB := First(strings, func(s string) bool {
		return s[0] == 'b'
	})

	if firstWithB != "bar" {
		t.Errorf("Expected first string starting with 'b' to be 'bar', got %s", firstWithB)
	}
}

func TestMap(t *testing.T) {
	// Test mapping numbers to their squares
	numbers := []int{1, 2, 3, 4}
	squares := Map(numbers, func(n int) int {
		return n * n
	})

	expected := []int{1, 4, 9, 16}
	if len(squares) != len(expected) {
		t.Errorf("Expected %d squares, got %d", len(expected), len(squares))
	}

	for i, square := range squares {
		if square != expected[i] {
			t.Errorf("Expected square of %d to be %d, got %d", numbers[i], expected[i], square)
		}
	}

	// Test mapping strings to their lengths
	strings := []string{"a", "ab", "abc", "abcd"}
	lengths := Map(strings, func(s string) int {
		return len(s)
	})

	expectedLengths := []int{1, 2, 3, 4}
	if len(lengths) != len(expectedLengths) {
		t.Errorf("Expected %d lengths, got %d", len(expectedLengths), len(lengths))
	}

	for i, length := range lengths {
		if length != expectedLengths[i] {
			t.Errorf("Expected length of %s to be %d, got %d", strings[i], expectedLengths[i], length)
		}
	}
}

func TestHas(t *testing.T) {
	// Test checking if collection has even number
	numbers := []int{1, 3, 5, 7, 9}
	hasEven := Has(numbers, func(n int) bool {
		return n%2 == 0
	})

	if hasEven {
		t.Errorf("Expected collection to not have even numbers")
	}

	numbers = append(numbers, 2)
	hasEven = Has(numbers, func(n int) bool {
		return n%2 == 0
	})

	if !hasEven {
		t.Errorf("Expected collection to have even numbers")
	}

	// Test checking if collection has string with specific prefix
	strings := []string{"foo", "bar", "baz"}
	hasPrefix := Has(strings, func(s string) bool {
		return len(s) > 0 && s[0] == 'q'
	})

	if hasPrefix {
		t.Errorf("Expected collection to not have strings starting with 'q'")
	}

	strings = append(strings, "qux")
	hasPrefix = Has(strings, func(s string) bool {
		return len(s) > 0 && s[0] == 'q'
	})

	if !hasPrefix {
		t.Errorf("Expected collection to have strings starting with 'q'")
	}
}

func TestReduce(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ret := Reduce(numbers, func(acc int, v int) int { return acc * v }, 1)
	if ret != 3628800 {
		t.Errorf("Expected 3628800, got %d", ret)
	}

	r := Reduce(numbers, func(l []int, v int) []int {
		return l
	}, nil)

	if len(r) == len(numbers) {
		t.Errorf("Expected %d, got %d", len(numbers), len(r))
	}
}
