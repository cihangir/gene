package gene

import (
	"os"

	"code.google.com/p/go.tools/imports"
)

func writeFormattedFile(fileName string, model []byte) error {

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
