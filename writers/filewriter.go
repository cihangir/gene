//Package writers handles the file write operations
package writers

import (
	"bytes"
	"go/format"
	"os"
	"regexp"

	"golang.org/x/tools/imports"
)

// WriteFormattedFile formats the code with goimports and writes the result to
// the given file, if file doesnt exists, it creates it
func WriteFormattedFile(fileName string, model []byte) error {
	dest, err := imports.Process("", model, nil)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(dest); err != nil {
		return err
	}

	return nil
}

// NewLinesRegex holds the regex to remove newlines from given bytes.Buffer
var NewLinesRegex = regexp.MustCompile(`(?m:\s*$)`)

// Clear formats the given source with predefined operations, it removes the
// new lines too
func Clear(buf bytes.Buffer) ([]byte, error) {
	bytes := NewLinesRegex.ReplaceAll(buf.Bytes(), []byte(""))

	// Format sources
	clean, err := format.Source(bytes)
	if err != nil {
		return buf.Bytes(), err
	}

	return clean, nil
}
