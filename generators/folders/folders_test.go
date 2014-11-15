package folders

import "testing"

func TestEnsureFolders(t *testing.T) {
	err := EnsureFolders("/tmp/", FolderStucture)
	if err != nil {
		t.Fatalf("err while creating folders: %s", err.Error())
	}

	// make sure trying to re-create doesnt give error
	err = EnsureFolders("/tmp/", FolderStucture)
	if err != nil {
		t.Fatalf("err while ensuring folders: %s", err.Error())
	}

	// // test relative path
	// err = EnsureFolders("./", folderStucture)
	// if err != nil {
	// 	t.Fatalf("err while writing to relative path: %s", err.Error())
	// }
}
