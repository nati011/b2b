package main

import (
	"flag"
	"fmt"
	"os"

	"marketplace/internal/tools/resourcegen"
)

func main() {
	var outputPath string
	flag.StringVar(&outputPath, "output", "config/resources.yaml", "Path to output resources manifest file")
	flag.StringVar(&outputPath, "o", "config/resources.yaml", "Path to output resources manifest file (short)")
	flag.Parse()

	if err := resourcegen.Generate(outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating resources manifest: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Resources manifest generated successfully at %s\n", outputPath)
}
