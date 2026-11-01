package main

import (
	"fmt"
	"strings"
)

func greet(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "Hello, World!"
	}
	return fmt.Sprintf("Hello, %s!", trimmed)
}
