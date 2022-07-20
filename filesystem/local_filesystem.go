package filesystem

import (
	"fmt"
	"os"
)

type LocalFilesystem struct{}

func NewLocalFilesystem() LocalFilesystem {
	return LocalFilesystem{}
}

func (fs LocalFilesystem) GetFile(path string) (Result, error) {
	buf, err := os.ReadFile(path)
	content := string(buf)
	if err != nil {
		fmt.Printf("%v\n", err, err)
		return nil, err
	}
	return FileResult{
		Path:    path,
		Content: content,
	}, nil
}

func (fs LocalFilesystem) GetFolder(path string) (Result, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var files []os.FileInfo
	for _, entry := range dirEntries {
		info, err := entry.Info()
		// TODO: should this fail the whole request, or return an "ChildError", i.e. "Unable to read"?
		// To show the client that the file exists, but is not "stat-able"
		if err != nil {
			return nil, err
		}
		// TODO: Do the requirements indicate this should ONLY return file results?
		if !info.IsDir() {
			files = append(files, info)
		}
	}
	children := make([]ChildResult, len(files))
	for i, file := range files {
		children[i] = ChildResult{
			Key:  file.Name(),
			Size: file.Size(),
		}
	}
	return FolderResult{
		Path:     path,
		Children: children,
	}, nil
}
