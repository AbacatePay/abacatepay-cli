// Package mock...
package mock

import (
	"math/rand"
	"strings"
)

func generateValidCPF(r *rand.Rand) string {
	digits := make([]int, 11)
	for i := range 9 {
		digits[i] = r.Intn(10)
	}

	sum := 0
	for i := range 9 {
		sum += digits[i] * (10 - i)
	}
	digits[9] = calculateDigit(sum)

	sum = 0
	for i := range 10 {
		sum += digits[i] * (11 - i)
	}
	digits[10] = calculateDigit(sum)

	var b strings.Builder

	for _, d := range digits {
		b.WriteByte('0' + byte(d))
	}

	return b.String()
}

func calculateDigit(sum int) int {
	remainder := (sum * 10) % 11
	if remainder < 10 {
		return remainder
	}
	return 0
}
