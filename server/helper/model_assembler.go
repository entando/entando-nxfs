package helper

import (
	"fmt"
	"github.com/entando/entando-nxfs/server/model"
	"os"
	"path"
	"path/filepath"
)

// ToDirectoryObject - create and return a DirectoryObject starting by the received path and FileInfo
func ToDirectoryObject(path string, fileInfo os.FileInfo) model.DirectoryObject {

	if fileInfo == nil {
		return model.DirectoryObject{}
	}

	objectType := model.F
	if fileInfo.IsDir() {
		objectType = model.D
	}

	rel, err := filepath.Rel(GetBrowsableFsRootPath(), path)
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
func ToDirectoryObjectFromFilePath(filePath string) model.DirectoryObject {

	var fileInfo os.FileInfo
	var err error

	if fileInfo, err = os.Stat(filePath); os.IsNotExist(err) {
		panic(err.Error())
	}

	fullPath := path.Dir(filePath)

	return ToDirectoryObject(fullPath, fileInfo)
}

// ToFileObjectFromFilePath - create and return a FileObject starting by the received path
func ToFileObjectFromFilePath(filePath string, content string) model.FileObject {

	var fileInfo os.FileInfo
	var err error

	if fileInfo, err = os.Stat(filePath); os.IsNotExist(err) {
		panic(err.Error())
	}

	fullPath := path.Dir(filePath)

	return ToFileObject(fullPath, fileInfo, content)
}

// ToFileObject - create and return a FileObject starting by the received path and FileInfo
func ToFileObject(path string, fileInfo os.FileInfo, fileContentString string) model.FileObject {

	if fileInfo == nil {
		return model.FileObject{}
	}

	rel, err := filepath.Rel(GetBrowsableFsRootPath(), path)
	if err != nil {
		panic("Error during relativization of the file " + fileInfo.Name())
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
