package nxsiteman

import "os"

// toDirectoryObject create and return a DirectoryObject starting by the received path and FileInfo
func toDirectoryObject(path string, fileInfo os.FileInfo) DirectoryObject {

	if nil == fileInfo {
		return DirectoryObject{}
	}

	objectType := F
	if fileInfo.IsDir() {
		objectType = D
	}

	return DirectoryObject{
		Name:    fileInfo.Name(),
		Path:    path,
		Size:    fileInfo.Size(),
		Type:    objectType,
		Created: ActionLog{At: fileInfo.ModTime()},
		Updated: ActionLog{At: fileInfo.ModTime()},
	}
}
