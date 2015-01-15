// Package folders creates required folders for module system
package folders

import (
	"os"
	"path"
)

// FolderStucture holds all the required paths to be created while setting up a
// new module
var FolderStucture = []string{
	"workers",
	"workers/",
	"models",
	"tests",
	"app",
}

// EnsureFolders checks and creates the folders for the modules
func EnsureFolders(root string, folderStucture []string) error {
	for _, folder := range folderStucture {
		path := path.Join(root, folder)
		exists, err := exists(path)
		if err != nil {
			return err
		}

		if exists {
			continue
		}

		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// check if folder or file exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
