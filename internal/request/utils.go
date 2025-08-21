package request

import "unicode"

func isAllCapitalLetters(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if !unicode.IsUpper(r) {
			return false
		}
	}

	return true
}
