package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddTableDriven(t *testing.T) {
	cases := map[string]struct {
		val1, val2 int
		expected   int
	}{
		"both_positive":      {2, 3, 5},
		"positive_plus_zero": {5, 0, 5},
		"negative_plus_pos":  {-1, 4, 3},
		"both_negative":      {-2, -3, -5},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := Add(tc.val1, tc.val2)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestSubtractTableDriven(t *testing.T) {
	scenarios := map[string]struct {
		v1, v2 int
		result int
	}{
		"both_positive":           {5, 2, 3},
		"positive_minus_zero":     {5, 0, 5},
		"negative_minus_positive": {-2, 3, -5},
		"both_negative":           {-5, -2, -3},
	}

	for name, s := range scenarios {
		t.Run(name, func(t *testing.T) {
			got := Subtract(s.v1, s.v2)
			require.Equal(t, s.result, got)
		})
	}
}

func TestDivide(t *testing.T) {
	t.Run("successful division", func(t *testing.T) {
		res, err := Divide(20, 4)
		require.NoError(t, err)
		require.Equal(t, 5, res)
	})

	t.Run("handling division by zero", func(t *testing.T) {
		_, err := Divide(10, 0)
		require.Error(t, err)
		require.EqualError(t, err, "division by zero")
	})
}