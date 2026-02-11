package cmd

import "fmt"

// Arch implements the targetArchitecture custom driver command.
// Apple Container only runs on Apple Silicon, so this always returns arm64.
func Arch() error {
	fmt.Print("arm64")
	return nil
}
