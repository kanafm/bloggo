package core

// Processor that takes a file and outputs a file
type Processor interface {
	// CanHandle determines whether this processor can handle this file
	CanHandle(inputFile string) bool

	// Handle processes file at file location, writing to output file
	Handle(inputFile string, outputFile string) error
}
