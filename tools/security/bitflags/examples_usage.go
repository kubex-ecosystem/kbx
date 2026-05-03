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

import gl "github.com/kubex-ecosystem/logz"

// Example of composing middleware decisions from flags.

// ExampleMwDecision shows how to compose middleware decisions from flags.
func ExampleMwDecision() {
	var sec SecFlag
	var reg FlagReg32[SecFlag]
	reg.Store(0)
	reg.Set(SecAuth | SecSanitize)
	sec = reg.Load()

	chain := []string{"trace", "logging"}
	if sec&SecAuth != 0 {
		chain = append(chain, "authentication")
	}
	if sec&SecSanitize != 0 {
		chain = append(chain, "sanitize")
	}
	gl.Println(chain)
	// Output: [trace logging authentication sanitize]
}
