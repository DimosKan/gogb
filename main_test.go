package main

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	var tests = []struct {
		a, b   int
		accept int
	}{
		{3, 3, 6},
		{2, 3, 5},
		{1, 5, 6},
		{-1, -3, -4},
		{-2, 5, 3},
	}
	for _, tt := range tests {
		results := fmt.Sprintf("%d,%d", tt.a, tt.b)
		t.Run(results, func(t *testing.T) {
			res := Add(tt.a, tt.b)
			if res != tt.accept {
				t.Errorf("Result is %d, we want %d", res, tt.accept)
			}
		})
	}
}
