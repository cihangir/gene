package writers

import (
	"bytes"
	"go/format"
	"os"
	"regexp"

	"golang.org/x/tools/imports"
)

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

	_, err = f.Write(dest)
	if err != nil {
		return err
	}

	return nil
}

var NewLinesRegex = regexp.MustCompile(`(?m:\s*$)`)

func Clear(buf bytes.Buffer) ([]byte, error) {
	bytes := NewLinesRegex.ReplaceAll(buf.Bytes(), []byte(""))

	// Format sources
	clean, err := format.Source(bytes)
	if err != nil {
		return buf.Bytes(), err
	}

	return clean, nil
}
