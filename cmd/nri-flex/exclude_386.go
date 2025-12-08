//go:build 386

package main

// BUILD ERROR: nri-flex no longer supports 32-bit (386) architecture.
//
// Windows 32-bit support has been removed from nri-flex.
// Please build with GOARCH=amd64 or GOARCH=arm64 instead.
//
// If you see this error, you are attempting to build for an unsupported architecture.

// Intentional compile error to prevent 386 builds
type _ struct {
	nri_flex_does_not_support_32bit_386_architecture int
}

var _ = nri_flex_does_not_support_32bit_386_architecture(0)
