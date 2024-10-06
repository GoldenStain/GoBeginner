package main

import (
	"fmt"
)

func SumInts(m map[string]int64) int64 {
	var res int64
	for _, v := range m {
		res += v
	}
	return res
}

func SumFloats(m map[string]float64) float64 {
	var res float64
	for _, v := range m {
		res += v
	}
	return res
}

func SumIntOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var res V
	for _, v := range m {
		res += v
	}
	return res
}

type Numbers interface {
	int64 | float64
}

func SumNumbers[K comparable, V Numbers](m map[K]V) V {
	var res V
	for _, v := range m {
		res += v
	}
	return res
}

func main() {
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	fmt.Printf("Non-generic sums: %v and %v\n",
		SumInts(ints),
		SumFloats(floats))

	fmt.Printf("Generic sums: %v and %v\n",
		SumIntOrFloats[string, int64](ints),
		SumIntOrFloats[string, float64](floats))

	fmt.Printf("Generic sums with Constraint: %v and %v\n",
		SumNumbers(ints),
		SumNumbers(floats))
}
