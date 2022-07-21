package filesystem

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
)

const testdataDirPath = "../testdata/local_filesystem/"

func TestLocalFilesystem_GetFolder(t *testing.T) {
	fs := NewLocalFilesystem()
	result, _ := fs.GetFolder(testdataDirPath)
	data := result.GetData()
	children := data["children"].([]ChildResult)
	if len(children) != 2 {
		t.Fail()
	}
	var childNames []string
	for _, child := range children {
		childNames = append(childNames, child.Key)
	}
	if reflect.DeepEqual(childNames, []string{"kv.txt", "data.json"}) {
		t.Fail()
	}
}

func TestLocalFilesystem_GetFolder_notexist(t *testing.T) {
	fs := NewLocalFilesystem()
	_, err := fs.GetFolder(fmt.Sprintf("%snot-exist/", testdataDirPath))
	if !errors.Is(err, os.ErrNotExist) {
		t.Fail()
	}
}

func TestLocalFilesystem_GetFile(t *testing.T) {
	fs := NewLocalFilesystem()
	result, _ := fs.GetFile(fmt.Sprintf("%skv.txt", testdataDirPath))
	data := result.GetData()
	content := data["content"].(string)
	expected := "ok=1"
	if content != expected {
		t.Fail()
	}
}

func TestLocalFilesystem_GetFile_notexist(t *testing.T) {
	fs := NewLocalFilesystem()
	_, err := fs.GetFile(fmt.Sprintf("%sdoes-not-exist.txt", testdataDirPath))
	if !errors.Is(err, os.ErrNotExist) {
		t.Fail()
	}
}
