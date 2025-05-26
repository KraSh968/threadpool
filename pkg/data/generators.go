package data

import "math/rand"

func GenerateInts(size int) *[]int {
	result := make([]int, size)
	for i := 0; i < len(result); i++ {
		result[i] = rand.Intn(100)
	}
	return &result
}

func GenerateRange(size int) *[]int {
	result := make([]int, size)
	for i := 0; i < len(result); i++ {
		result[i] = i
	}
	return &result
}

func MultiplyInts(ints *[]int, multiplier int) *[]int {
	for i := 0; i < len(*ints); i++ {
		(*ints)[i] = (*ints)[i] * multiplier
	}
	return ints
}
