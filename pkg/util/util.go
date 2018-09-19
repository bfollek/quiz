// Package util has utility code that doesn't fit into any other package.
package util

import (
	"bufio"
	"fmt"
)

// GetInput prompts and reads input.
func GetInput(r *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	return r.ReadString('\n')
}
