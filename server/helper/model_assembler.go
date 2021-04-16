package helper

import (
	"fmt"
	"github.com/entando/entando-nxfs/server/model"
	"os"
	"path"
	"path/filepath"
)

type BrowsableFsRootPathGetter func() string

type ModelAssembler struct {
	browsableFsRootPathGetter BrowsableFsRootPathGetter
}

// NewModelAssembler creates a ModelAssembler
func NewModelAssembler(browsableFsRootPathGetter BrowsableFsRootPathGetter) ModelAssembler {
	return ModelAssembler{browsableFsRootPathGetter}
}

// ToDirectoryObject - create and return a DirectoryObject starting by the received path and FileInfo
func (m *ModelAssembler) ToDirectoryObject(path string, fileInfo os.FileInfo) model.DirectoryObject {

	if fileInfo == nil {
		return model.DirectoryObject{}
	}

	objectType := model.F
	if fileInfo.IsDir() {
		objectType = model.D
	}

	rel, err := filepath.Rel(m.browsableFsRootPathGetter(), path)
	if err != nil {
		panic("Error during relativization of the file " + fileInfo.Name())
	}

	return model.DirectoryObject{
		Name:    fileInfo.Name(),
		Path:    rel,
		Size:    fileInfo.Size(),
		Type:    objectType,
		Created: model.ActionLog{At: fileInfo.ModTime()},
		Updated: model.ActionLog{At: fileInfo.ModTime()},
	}
}

// ToDirectoryObjectFromFilePath - create and return a DirectoryObject starting by the received path
func (m *ModelAssembler) ToDirectoryObjectFromFilePath(filePath string) model.DirectoryObject {

	var fileInfo os.FileInfo
	var err error

	if fileInfo, err = os.Stat(filePath); os.IsNotExist(err) {
		return model.DirectoryObject{}
	}

	fullPath := path.Dir(filePath)

	return m.ToDirectoryObject(fullPath, fileInfo)
}

// ToFileObjectFromFilePath - create and return a FileObject starting by the received path
func (m *ModelAssembler) ToFileObjectFromFilePath(filePath string, content string) model.FileObject {

	var fileInfo os.FileInfo
	var err error

	if fileInfo, err = os.Stat(filePath); os.IsNotExist(err) {
		return model.FileObject{}
	}

	fullPath := path.Dir(filePath)

	return m.ToFileObject(fullPath, fileInfo, content)
}

// ToFileObject - create and return a FileObject starting by the received path and FileInfo
func (m *ModelAssembler) ToFileObject(path string, fileInfo os.FileInfo, fileContentString string) model.FileObject {

	if fileInfo == nil {
		return model.FileObject{}
	}

	rel, err := filepath.Rel(m.browsableFsRootPathGetter(), path)
	if err != nil {
		return model.FileObject{}
	}

	return model.FileObject{
		Name:    fileInfo.Name(),
		Path:    rel,
		Size:    fileInfo.Size(),
		Type:    model.F,
		Created: model.ActionLog{At: fileInfo.ModTime()},
		Updated: model.ActionLog{At: fileInfo.ModTime()},
		Content: fmt.Sprintf(fileContentString),
	}
}
