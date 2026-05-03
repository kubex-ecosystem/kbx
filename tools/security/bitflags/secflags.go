// Package bitflags provides atomic flag registers and utilities for managing
// sets of bit flags on top of basic integer types.
//
// The package defines generic flag types `FlagReg32[T ~uint32]` and
// `FlagReg64[T ~uint64]`, along with helper utilities. The type parameter `T`
// must be an unsigned integer type.
//
// Operations like Set, Clear, Toggle, and Mask are implemented using atomic
// compare-and-swap (CAS) loops, ensuring thread safety without requiring
// external mutexes.
//
// Conditional helpers `SetIf` and `ClearIf` allow modifying flags only when
// certain other flags are set or clear respectively, with atomic guarantees.
//
// Pretty-print utilities provide formatted string representations of flag
// combinations, mapping flag values to symbolic names.
//
// Example:
//
//	type Perms uint32
//
//	const (
//	    Read  Perms = 1 << iota
//	    Write
//	    Exec
//	)
//
//	var p FlagReg32[Perms]
//	p.Set(Read | Write)
//	if p.Has(Read) { fmt.Println("readable") }
package bitflags

// SecFlag represents a set of security flags.
type SecFlag uint32

const (
	// SecAuth enables authentication.
	SecAuth SecFlag = 1 << iota
	// SecSanitize enables sanitization.
	SecSanitize
	// SecSanitizeBody enables body sanitization.
	SecSanitizeBody
)

// FromLegacyMap bridges your map[string]bool → flags.
func FromLegacyMap(m map[string]bool) SecFlag {
	var f SecFlag
	if m["secure"] {
		f |= SecAuth
	}
	if m["validateAndSanitize"] {
		f |= SecSanitize
	}
	if m["validateAndSanitizeBody"] {
		f |= SecSanitizeBody
	}
	return f
}

var secNames = map[string]SecFlag{
	"auth":          SecAuth,
	"sanitize":      SecSanitize,
	"sanitize_body": SecSanitizeBody,
}

func (f SecFlag) String() string { return FlagString(f, secNames) }
