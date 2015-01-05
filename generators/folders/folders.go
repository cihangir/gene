package folders

import (
	"os"
	"path"
)

var FolderStucture = []string{
	"workers",
	"workers/",
	"models",
	"tests",
	"app",
}

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
