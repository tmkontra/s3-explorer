package filesystem

const (
	FileResultType   = "object"
	FolderResultType = "folder"
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
	GetPath(path string) (Result, error)
}
