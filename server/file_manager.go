package nxsiteman

import (
	"fmt"
	pkgErr "github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// isDirWithChildren
func isDirWithChildren(absPathFile string, fileToDelete os.FileInfo) bool {
	if children, _ := ioutil.ReadDir(absPathFile); fileToDelete.IsDir() && len(children) > 0 {
		return true
	} else {
		return false
	}
}

// getFileInfoIfPathExistOrErrorResponse - if the received path exists return the corresponding os.FileInfo, otherwise return an error ImplResponse
func getFileInfoIfPathExistOrErrorResponse(pathToCheck string) (os.FileInfo, *ImplResponse) {
	if fileInfo, err := os.Stat(pathToCheck); os.IsNotExist(err) {
		return nil, ErrorResponse(http.StatusNotFound, "path_not_found", err.Error())
	} else {
		return fileInfo, nil
	}
}

// getDraftPageInfoIfExistOrErrorResponse - if the received path exists as relative path in the draft pages folder return the corresponding os.FileInfo, otherwise return an error ImplResponse
func getDraftPageInfoIfExistOrErrorResponse(pathToCheck string) (os.FileInfo, *ImplResponse) {
	return getFileInfoIfPathExistOrErrorResponse(path.Join(GetDraftPagesPath(), pathToCheck))
}

// getPublishedPageInfoIfExistOrErrorResponse - if the received path exists as relative path in the published pages folder return the corresponding os.FileInfo, otherwise return an error ImplResponse
func getPublishedPageInfoIfExistOrErrorResponse(pathToCheck string) (os.FileInfo, *ImplResponse) {
	return getFileInfoIfPathExistOrErrorResponse(path.Join(GetPublishedPagesPath(), pathToCheck))
}

// relativizeToDraftPageFolder - receive a path and relativize it to the draft pages folder
func relativizeToDraftPageFolder(pathToRelativize string) string {
	return path.Join(GetDraftPagesPath(), pathToRelativize)
}

// relativizeToPublishedPageFolder - receive a path and relativize it to the draft pages folder
func relativizeToPublishedPageFolder(pathToRelativize string) string {
	return path.Join(GetPublishedPagesPath(), pathToRelativize)
}

// createFile - create a file in the received path containing the received content return an error ImplResponse if an error occurs, nil otherwise
func createFile(path string, fileObject FileObject) (errorResp *ImplResponse) {

	data := []byte(fileObject.Content)
	err := ioutil.WriteFile(path, data, 0755)
	if err != nil {
		errorResp = ErrorResponse(http.StatusBadRequest, "write_error", err.Error())
	}

	return errorResp
}

// createDirectory - create a directory in the received path. return an error ImplResponse if an error occurs, nil otherwise
func createDirectory(dirPath string) (errorResp *ImplResponse) {

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			errorResp = ErrorResponse(http.StatusBadRequest, "dir_write_error", err.Error())
		}
	}

	return errorResp
}

// decodePath - receives an url encoded path, decodes and returns it. if a decode error occurs, it will return an error ImplResponse
func decodePath(encodedPath string) (string, *ImplResponse) {

	decodedPath, err := url.PathUnescape(encodedPath)
	if err != nil {
		return "", ErrorResponse(http.StatusBadRequest, "error_decoding_path", err.Error())
	}

	return decodedPath, nil
}

// deleteFile - delete a file or folder, return an error ImplResponse if an error occur, nil otherwise
func deleteFile(filePath string) *ImplResponse {
	err := os.Remove(filePath)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, "deletion_error", "An error occurred during the deletion")
	} else {
		return nil
	}
}

// browseFileTree - traverse recursively the path represented by fileInfo
func browseFileTree(path string, fileInfo os.FileInfo, currDepth int32, maxDepth int32, directoryObjects []DirectoryObject) ([]DirectoryObject, error) {

	// if depth reached return
	if currDepth > maxDepth && maxDepth != 0 {
		return directoryObjects, nil
	}

	// if the current one is a file add it to the result list and return
	if !fileInfo.IsDir() {
		directoryObjects = append(directoryObjects, toDirectoryObject(path, fileInfo))
		return directoryObjects, nil
	}

	// otherwise proceed with the tree inspection
	dirAbsPath := filepath.Join(path, fileInfo.Name())

	// read dir
	readFilesInfo, err := ioutil.ReadDir(dirAbsPath)
	if err != nil {
		return directoryObjects, pkgErr.Wrap(err, fmt.Sprintf("can't read directory %s", dirAbsPath))
	}

	// call recursively
	for _, file := range readFilesInfo {
		directoryObjects, err = browseFileTree(dirAbsPath, file, currDepth+1, maxDepth, directoryObjects)
		if err != nil {
			return directoryObjects, err
		}
	}

	return directoryObjects, nil
}

// composeFullPathOrErrorResponse - receives a URL encoded path, decodes it and return the corresponding full path (without filename), the fileInfo of the requested file/folder and a possible REST response containing an error
func composeFullPathOrErrorResponse(encodedPath string) (fullPath string, fileInfoToBrowse os.FileInfo, errorResponse *ImplResponse) {

	// define paths
	decodedPath, errResponse := decodePath(encodedPath)
	if errResponse != nil {
		return "", nil, errResponse
	}

	fullPathToBrowse := filepath.Join(GetBrowsableFsRootPath(), decodedPath)

	// does path exist?
	fileInfoToBrowse, errRespFileExist := getFileInfoIfPathExistOrErrorResponse(fullPathToBrowse)
	if nil != errRespFileExist {
		return "", nil, errRespFileExist
	}

	// extract pathToBrowse path
	fullPath = path.Dir(fullPathToBrowse)

	return fullPath, fileInfoToBrowse, nil
}

// copyFileTo - copy the file identified by originFile to the destination identified by destinationFile. return an error ImplResponse if an error occur, nil otherwise
func copyFileTo(sourceFile string, destinationFile string) *ImplResponse {

	// create path if needed
	destinationFolder := path.Dir(destinationFile)
	if _, implResponse := getFileInfoIfPathExistOrErrorResponse(destinationFolder); nil != implResponse {
		err := os.MkdirAll(destinationFolder, os.ModePerm)
		if err != nil {
			return ErrorResponse(http.StatusInternalServerError, "path_creation_error", "An error occurred during the creation of the published page path")
		}
	}

	// input
	in, err := os.Open(sourceFile)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, "draft_read_error", fmt.Sprintf("An error occurred during the read operation of the draft page file %s", sourceFile))
	}
	defer in.Close()

	// output
	out, err := os.Create(destinationFile)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, "published_creation_error", fmt.Sprintf("An error occurred during the creation of the published page file %s", destinationFile))
	}
	defer out.Close()

	// copy file
	_, err = io.Copy(out, in)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, "published_copy_error", "An error occurred during the copy of the draft page file to the the published page file")
	} else {
		return nil
	}
}
