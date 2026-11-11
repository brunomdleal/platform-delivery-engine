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
	blue   = color.New(color.FgBlue, color.Bold).SprintfFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

var spacing = 21

// Info Prints logs with a cyan-badged
func Info(format string, a ...interface{}) {
	fmt.Printf("%*s %s\n", spacing, cyan("[INFO]"), fmt.Sprintf(format, a...))
}

// Success Prints logs with a green-badged
func Success(format string, a ...interface{}) {
	fmt.Printf("%*s %s\n", spacing, green("[✓]"), fmt.Sprintf(format, a...))
}

// Warn Prints logs with a yellow-badged
func Warn(format string, a ...interface{}) {
	fmt.Printf("%*s %s\n", spacing, yellow("[!]"), fmt.Sprintf(format, a...))
}

// Error Prints logs with a red-badged
func Error(format string, a ...interface{}) {
	fmt.Printf("%*s %s\n", spacing, red("[ERROR]"), fmt.Sprintf(format, a...))
}

// Header Prints a formatted section banner
func Header(title string) {
	fmt.Println()
	fmt.Printf("%*s %s\n", spacing, cyan("==>"), bold(title))
}

// Version Prints the tool version
func Version(title string) {
	fmt.Println()
	fmt.Printf("%s %s\n", "pde", bold(title))
}

// Directory Prints the directory with a blue-badged
func Directory(format string, a ...interface{}) {
	fmt.Printf("%*s %s\n", spacing, blue("[DIR]"), fmt.Sprintf(format, a...))
}
