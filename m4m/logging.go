package m4m

import (
	"fmt"
	"os"
)

// LogError logs into stderr
func LogError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format + "\n", a...)
}

// LogInfo logs into stdout
func LogInfo(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, format + "\n", a...)
}
