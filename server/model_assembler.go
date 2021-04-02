package nxsiteman

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// toDirectoryObject - create and return a DirectoryObject starting by the received path and FileInfo
func toDirectoryObject(path string, fileInfo os.FileInfo) DirectoryObject {

	if nil == fileInfo {
		return DirectoryObject{}
	}

	objectType := F
	if fileInfo.IsDir() {
		objectType = D
	}

	rel, err := filepath.Rel(GetBrowsableFsRootPath(), path)
	if nil != err {
		panic("Error during relativization of the file " + fileInfo.Name())
	}

	return DirectoryObject{
		Name:    fileInfo.Name(),
		Path:    rel,
		Size:    fileInfo.Size(),
		Type:    objectType,
		Created: ActionLog{At: fileInfo.ModTime()},
		Updated: ActionLog{At: fileInfo.ModTime()},
	}
}

// toDirectoryObject - create and return a DirectoryObject starting by the received path
func toDirectoryObjectFromFilePath(filePath string) DirectoryObject {

	var fileInfo os.FileInfo
	var err error

	if fileInfo, err = os.Stat(filePath); os.IsNotExist(err) {
		panic(err.Error())
	}

	fullPath := path.Dir(filePath)

	return toDirectoryObject(fullPath, fileInfo)
}

// toFileObject - create and return a FileObject starting by the received path and FileInfo
func toFileObject(path string, fileInfo os.FileInfo, fileContentString string) FileObject {

	if nil == fileInfo {
		return FileObject{}
	}

	rel, err := filepath.Rel(GetBrowsableFsRootPath(), path)
	if nil != err {
		panic("Error during relativization of the file " + fileInfo.Name())
	}

	return FileObject{
		Name:    fileInfo.Name(),
		Path:    rel,
		Size:    fileInfo.Size(),
		Type:    F,
		Created: ActionLog{At: fileInfo.ModTime()},
		Updated: ActionLog{At: fileInfo.ModTime()},
		Content: fmt.Sprintf(fileContentString),
	}
}
