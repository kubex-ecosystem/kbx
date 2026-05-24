package is

import (
	"regexp"
	"unicode"
)

// Base64String checks if a string is a valid Base64 encoded string.
func Base64String(s string) bool {
	matched, _ := regexp.MatchString("^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$", s)
	return matched
}

// Base64ByteSlice checks if a byte slice is a valid Base64 encoded slice.
func Base64ByteSlice(s []byte) bool {
	matched, _ := regexp.Match("^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$", s)
	return matched
}

// Base64ByteSliceString checks if the byte representation of a string is a valid Base64 encoded slice.
func Base64ByteSliceString(s string) bool {
	matched, _ := regexp.Match("^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$", []byte(s))
	return matched
}

// Base64ByteSliceStringWithPadding checks if the byte representation of a string is a valid Base64 encoded slice (with padding support).
func Base64ByteSliceStringWithPadding(s string) bool {
	matched, _ := regexp.Match("^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$", []byte(s))
	return matched
}

// URLEncodeString checks if a string consists entirely of URL-safe characters.
func URLEncodeString(s string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9%_.-]+$", s)
	return matched
}

// URLEncodeByteSlice checks if a byte slice consists entirely of URL-safe characters.
func URLEncodeByteSlice(s []byte) bool {
	matched, _ := regexp.Match("^[a-zA-Z0-9%_.-]+$", s)
	return matched
}

// Base62String checks if a string is a valid Base62 encoded string.
// It returns false if the string starts with a digit.
func Base62String(s string) bool {
	if unicode.IsDigit(rune(s[0])) {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", s)
	return matched
}

// Base62ByteSlice checks if a byte slice is a valid Base62 encoded byte slice.
func Base62ByteSlice(s []byte) bool {
	matched, _ := regexp.Match("^[a-zA-Z0-9_]+$", s)
	return matched
}
