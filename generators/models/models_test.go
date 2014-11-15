package models

import (
	"encoding/json"

	"testing"

	"bitbucket.org/cihangirsavas/gene/generators/folders"
	"bitbucket.org/cihangirsavas/gene/schema"
	"bitbucket.org/cihangirsavas/gene/testdata"
	"bitbucket.org/cihangirsavas/gene/writers"
)

func TestGenerateModel(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	model, err := s.GenerateModel()
	if err != nil {
		t.Fatal(err.Error())
	}

	folders.EnsureFolders("/tmp/", folders.FolderStucture)
	fileName := "/tmp/gene/models/" + s.Title + ".go"

	err = writers.WriteFormattedFile(fileName, model)
	if err != nil {
		t.Fatal(err.Error())
	}

}
