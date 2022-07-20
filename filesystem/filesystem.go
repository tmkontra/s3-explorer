package filesystem

import "strings"

const (
	FileResultType        = "object"
	FolderResultType      = "folder"
	folderSuffixCharacter = "/" // Some filesystems use a different directory separator
)

type Result interface {
	GetData() map[string]interface{}
}

type FileResult struct {
	Path    string
	Content string
}

func (f FileResult) GetData() map[string]interface{} {
	return map[string]interface{}{
		"type":    FileResultType,
		"content": f.Content,
	}
}

type ChildResult struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}

type FolderResult struct {
	Path     string
	Children []ChildResult
}

func (f FolderResult) GetData() map[string]interface{} {
	return map[string]interface{}{
		"type":     FolderResultType,
		"children": f.Children,
	}
}

type Filesystem interface {
	GetFolder(path string) (Result, error)
	GetFile(path string) (Result, error)
}

func GetPath(fs Filesystem, path string) (Result, error) {
	if strings.HasSuffix(path, folderSuffixCharacter) {
		return fs.GetFolder(path)
	} else {
		return fs.GetFile(path)
	}
}
