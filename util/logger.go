package util

import "fmt"

func SetErrorDetails(cause, effect, solution, position string) string {
	return fmt.Sprintf("[Cause] %s\n[Effect] %s\n[Solution] %s\n[Position] %s:", cause, effect, solution, position)
}
