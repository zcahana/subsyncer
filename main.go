package main

import (
	"flag"
)

// Flags
var (
	inputFile string
	inputLanguage string

	referenceFile string
	referenceLanguage string
)
func main() {
	flag.StringVar(&inputFile, "input-file", "", "Path to subtitle file to synchronize")
	flag.StringVar(&inputLanguage, "input-lang", "", "Language of subtitle file to synchronize")
	flag.StringVar(&referenceFile, "ref-file", "", "Path to reference subtitle file")
	flag.StringVar(&referenceLanguage, "ref-lang", "", "Langauge of reference subtitle file")

	flag.Parse()
}