package helper

import "strconv"

// AtoiSafe converts string to int safely.
// Returns 0 if conversion fails.
func AtoiSafe(s string) int {
	if s == "" {
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
