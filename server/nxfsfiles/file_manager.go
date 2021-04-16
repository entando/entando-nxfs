package nxfsfiles

import (
	"fmt"
	"github.com/entando/entando-nxfs/server/helper"
	"github.com/entando/entando-nxfs/server/model"
	"github.com/entando/entando-nxfs/server/net"
	pkgErr "github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

const DraftPagesDir = "pages"
const PublishedPagesDir = "pub_pages"

type FileManager struct {
	modelAssembler helper.ModelAssembler
}

// NewFileManager creates a FileManager
func NewFileManager(modelAssembler helper.ModelAssembler) FileManager {
	return FileManager{modelAssembler: modelAssembler}
}

// IsDirWithChildren - return true if path is a dir and has childre, false otherwise
func (fm *FileManager) IsDirWithChildren(absPathFile string, fileToDelete os.FileInfo) bool {
	if children, _ := ioutil.ReadDir(absPathFile); fileToDelete.IsDir() && len(children) > 0 {
		return true
	} else {
		return false
	}
}

// GetFileInfoIfPathExistOrErrorResponse - if the received path exists return the corresponding os.FileInfo, otherwise return an error NxfsResponse
func (fm *FileManager) GetFileInfoIfPathExistOrErrorResponse(pathToCheck string) (os.FileInfo, *net.NxfsResponse) {
	if fileInfo, err := os.Stat(pathToCheck); os.IsNotExist(err) {
		return nil, helper.ErrorResponse(http.StatusNotFound, "path_not_found", err.Error())
	} else {
		return fileInfo, nil
	}
}

// GetDraftPageInfoIfExistOrErrorResponse - if the received path exists as relative path in the draft pages folder return the corresponding os.FileInfo, otherwise return an error NxfsResponse
func (fm *FileManager) GetDraftPageInfoIfExistOrErrorResponse(pathToCheck string) (os.FileInfo, *net.NxfsResponse) {
	return fm.GetFileInfoIfPathExistOrErrorResponse(path.Join(helper.GetDraftPagesPath(), pathToCheck))
}

// getPublishedPageInfoIfExistOrErrorResponse - if the received path exists as relative path in the published pages folder return the corresponding os.FileInfo, otherwise return an error NxfsResponse
func (fm *FileManager) getPublishedPageInfoIfExistOrErrorResponse(pathToCheck string) (os.FileInfo, *net.NxfsResponse) {
	return fm.GetFileInfoIfPathExistOrErrorResponse(path.Join(helper.GetPublishedPagesPath(), pathToCheck))
}

// RelativizeToDraftPageFolder - receive a path and relativize it to the draft pages folder
func (fm *FileManager) RelativizeToDraftPageFolder(pathToRelativize string) string {
	return path.Join(helper.GetDraftPagesPath(), pathToRelativize)
}

// RelativizeToPublishedPageFolder - receive a path and relativize it to the draft pages folder
func (fm *FileManager) RelativizeToPublishedPageFolder(pathToRelativize string) string {
	return path.Join(helper.GetPublishedPagesPath(), pathToRelativize)
}

// CreateFile - create a file in the received path containing the received content return an error NxfsResponse if an error occurs, nil otherwise
func (fm *FileManager) CreateFile(path string, fileObject model.FileObject) (errorResp *net.NxfsResponse) {

	data := []byte(fileObject.Content)
	err := ioutil.WriteFile(path, data, 0755)
	if err != nil {
		errorResp = helper.ErrorResponse(http.StatusBadRequest, "write_error", err.Error())
	}

	return errorResp
}

// CreateDirectory - create a directory in the received path. return an error NxfsResponse if an error occurs, nil otherwise
func (fm *FileManager) CreateDirectory(dirPath string) (errorResp *net.NxfsResponse) {

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			errorResp = helper.ErrorResponse(http.StatusBadRequest, "dir_write_error", err.Error())
		}
	}

	return errorResp
}

// DecodePath - receives an url encoded path, decodes and returns it. if a decode error occurs, it will return an error NxfsResponse
func (fm *FileManager) DecodePath(encodedPath string) (string, *net.NxfsResponse) {

	decodedPath, err := url.PathUnescape(encodedPath)
	if err != nil {
		return "", helper.ErrorResponse(http.StatusBadRequest, "error_decoding_path", err.Error())
	}

	return decodedPath, nil
}

// DeleteFile - delete a file or folder, return an error NxfsResponse if an error occur, nil otherwise
func (fm *FileManager) DeleteFile(filePath string) *net.NxfsResponse {
	err := os.Remove(filePath)
	if err != nil {
		return helper.ErrorResponse(http.StatusInternalServerError, "deletion_error", "An error occurred during the deletion")
	} else {
		return nil
	}
}

// BrowseFileTree - traverse recursively the path represented by fileInfo
func (fm *FileManager) BrowseFileTree(path string, fileInfo os.FileInfo, currDepth int32, maxDepth int32, directoryObjects []model.DirectoryObject, showPublishedPages bool) ([]model.DirectoryObject, error) {

	// if depth reached return
	if currDepth > maxDepth && maxDepth != 0 ||
		(showPublishedPages && fileInfo.Name() == DraftPagesDir) ||
		(!showPublishedPages && fileInfo.Name() == PublishedPagesDir) {
		return directoryObjects, nil
	}

	// if the current one is a file add it to the result list and return
	if !fileInfo.IsDir() {
		directoryObjects = append(directoryObjects, fm.modelAssembler.ToDirectoryObject(path, fileInfo))
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
		directoryObjects, err = fm.BrowseFileTree(dirAbsPath, file, currDepth+1, maxDepth, directoryObjects, showPublishedPages)
		if err != nil {
			return directoryObjects, err
		}
	}

	return directoryObjects, nil
}

// ComposeFullPathOrErrorResponse - receives a URL encoded path, decodes it and return the corresponding full path (without filename), the fileInfo of the requested file/folder and a possible REST response containing an error
func (fm *FileManager) ComposeFullPathOrErrorResponse(encodedPath string) (fullPath string, fileInfoToBrowse os.FileInfo, errorResponse *net.NxfsResponse) {

	// define paths
	decodedPath, errResponse := fm.DecodePath(encodedPath)
	if errResponse != nil {
		return "", nil, errResponse
	}

	fullPathToBrowse := filepath.Join(helper.GetBrowsableFsRootPath(), decodedPath)

	// does path exist?
	fileInfoToBrowse, errRespFileExist := fm.GetFileInfoIfPathExistOrErrorResponse(fullPathToBrowse)
	if nil != errRespFileExist {
		return "", nil, errRespFileExist
	}

	// extract pathToBrowse path
	fullPath = path.Dir(fullPathToBrowse)

	return fullPath, fileInfoToBrowse, nil
}

// CopyFileTo - copy the file identified by originFile to the destination identified by destinationFile. return an error NxfsResponse if an error occur, nil otherwise
func (fm *FileManager) CopyFileTo(sourceFile string, destinationFile string) *net.NxfsResponse {

	// create path if needed
	destinationFolder := path.Dir(destinationFile)
	if _, implResponse := fm.GetFileInfoIfPathExistOrErrorResponse(destinationFolder); nil != implResponse {
		err := os.MkdirAll(destinationFolder, os.ModePerm)
		if err != nil {
			return helper.ErrorResponse(http.StatusInternalServerError, "path_creation_error", "An error occurred during the creation of the published page path")
		}
	}

	// input
	in, err := os.Open(sourceFile)
	if err != nil {
		return helper.ErrorResponse(http.StatusInternalServerError, "draft_read_error", fmt.Sprintf("An error occurred during the read operation of the draft page file %s", sourceFile))
	}
	defer in.Close()

	// output
	out, err := os.Create(destinationFile)
	if err != nil {
		return helper.ErrorResponse(http.StatusInternalServerError, "published_creation_error", fmt.Sprintf("An error occurred during the creation of the published page file %s", destinationFile))
	}
	defer out.Close()

	// copy file
	_, err = io.Copy(out, in)
	if err != nil {
		return helper.ErrorResponse(http.StatusInternalServerError, "published_copy_error", "An error occurred during the copy of the draft page file to the the published page file")
	} else {
		return nil
	}
}
