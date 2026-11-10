package output

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	red    = color.New(color.FgRed, color.Bold).SprintFunc()
	cyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

// Info Prints logs with a cyan-badged
func Info(format string, a ...interface{}) {
	fmt.Printf("%s %s\n", cyan("[INFO]"), fmt.Sprintf(format, a...))
}

// Success Prints logs with a cyan-badged
func Success(format string, a ...interface{}) {
	fmt.Printf("%s %s\n", green("[✓]"), fmt.Sprintf(format, a...))
}

// Warn Prints logs with a cyan-badged
func Warn(format string, a ...interface{}) {
	fmt.Printf("%s %s\n", yellow("[!]"), fmt.Sprintf(format, a...))
}

// Error Prints logs with a cyan-badged
func Error(format string, a ...interface{}) {
	fmt.Printf("%s %s\n", yellow("[ERROR]"), fmt.Sprintf(format, a...))
}

// Header Prints a formatted section banner
func Header(title string) {
	fmt.Println()
	fmt.Printf("%s %s\n", cyan("==>"), bold(title))
}

// Version Prints the tool version
func Version(title string) {
	fmt.Println()
	fmt.Printf("%s %s\n", "pde", bold(title))
}
