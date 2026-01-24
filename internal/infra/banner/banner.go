package banner

import (
	"fmt"
	"strings"
)

// Print displays the ASCII art banner for NUCLEUS with version and environment info
func Print(version, env string) {
	banner := "       __                 _        \n" +
		"      / _|               | |       \n" +
		"  ___| |_ ___  _   _  ___| |_ __ _ \n" +
		" / _ \\  _/ _ \\| | | |/ _ \\ __/ _` |\n" +
		"|  __/ || (_) | |_| |  __/ || (_| |\n" +
		" \\___|_| \\___/ \\__, |\\___|\\__\\__,_|\n" +
		"                __/ |              \n" +
		"               |___/               "
	tagline := "\n Marketplace Platform"

	// Calculate banner width by finding the longest line
	bannerLines := strings.Split(strings.TrimRight(banner, "\n"), "\n")
	bannerWidth := 0
	for _, line := range bannerLines {
		if len(line) > bannerWidth {
			bannerWidth = len(line)
		}
	}

	separator := strings.Repeat("=", bannerWidth)

	fmt.Println()
	fmt.Println(separator)
	fmt.Println()

	fmt.Print(banner)

	// Center tagline relative to banner width
	centerText(tagline, bannerWidth)
	fmt.Println()

	fmt.Println(separator)
	fmt.Println()

	// Format info lines with aligned labels
	labelWidth := len("Environment:")
	infoLines := []struct {
		label string
		value string
	}{
		{"Version:", version},
		{"Environment:", strings.ToUpper(env)},
	}

	// Calculate left padding to center the block
	maxLineWidth := labelWidth + 1 + len(version) // "Version: " + version
	if len(strings.ToUpper(env)) > len(version) {
		maxLineWidth = labelWidth + 1 + len(strings.ToUpper(env))
	}
	leftPadding := (bannerWidth - maxLineWidth) / 2

	for _, info := range infoLines {
		// Align labels by padding to consistent position
		fmt.Print(strings.Repeat(" ", leftPadding))
		fmt.Printf("%-*s %s\n", labelWidth, info.label, info.value)
	}

	fmt.Println()
	fmt.Println(separator)
	fmt.Println()
}

// centerText centers the given text within the specified width
func centerText(text string, width int) {
	padding := (width - len(text)) / 2
	if padding > 0 {
		fmt.Print(strings.Repeat(" ", padding))
	}
	fmt.Print(text)
	if padding > 0 {
		fmt.Print(strings.Repeat(" ", width-len(text)-padding))
	}
}
