package data

import "math/rand"

func GenerateInts(size int) *[]int {
	result := make([]int, size)
	for i := 0; i < len(result); i++ {
		result[i] = rand.Intn(100)
	}
	return &result
}
