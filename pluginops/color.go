package main

import (
	"fmt"
)

// From https://twin.sh/articles/35/how-to-add-colors-to-your-console-terminal-output-in-go
var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"

func Red(s string) string {
	return fmt.Sprintf("%s%s%s", red, s, reset)
}

func Green(s string) string {
	return fmt.Sprintf("%s%s%s", green, s, reset)
}
