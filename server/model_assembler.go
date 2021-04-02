package nxsiteman

import (
	"os"
	"path/filepath"
)

// toDirectoryObject create and return a DirectoryObject starting by the received path and FileInfo
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
