package nxfspages

import (
	"github.com/entando/entando-nxfs/server/helper"
	"github.com/entando/entando-nxfs/server/net"
	"github.com/entando/entando-nxfs/server/nxfsfiles"
	"net/http"
	"os"
	"strings"
)

const pageSuffix = ".page"

// PublishPage - publish the received draft page and return an error NxfsResponse if an error occurs, nil otherwise
func PublishPage(encodedDraftPagePath string) (errorResp *net.NxfsResponse) {

	var pageFileInfo os.FileInfo

	// decode path
	decodedPath, errResponse := nxfsfiles.DecodePath(encodedDraftPagePath)
	if errResponse != nil {
		return errResponse
	}

	suffixedPage := addPageSuffix(decodedPath)

	// check if file exist as draft in the correct folder or error
	pageFileInfo, errResponse = nxfsfiles.GetDraftPageInfoIfExistOrErrorResponse(suffixedPage)
	if errResponse != nil {
		return errResponse
	}

	// if dir error
	if pageFileInfo.IsDir() {
		return helper.ErrorResponse(http.StatusUnprocessableEntity, "cannot_publish_dir", "The received path corresponds to a directory, only pages can be published")
	}

	draftPageFullPath := nxfsfiles.RelativizeToDraftPageFolder(suffixedPage)
	publishedPageFullPath := nxfsfiles.RelativizeToPublishedPageFolder(suffixedPage)
	return nxfsfiles.CopyFileTo(draftPageFullPath, publishedPageFullPath)
}

// UnpublishPage - unpublish the received published page and return an error NxfsResponse if an error occurs, nil otherwise
func UnpublishPage(encodedPublishedPagePath string) (errorResp *net.NxfsResponse) {

	// decode path
	decodedPath, errResponse := nxfsfiles.DecodePath(encodedPublishedPagePath)
	if errResponse != nil {
		return errResponse
	}

	suffixedPage := addPageSuffix(decodedPath)
	publishedPageFullPath := nxfsfiles.RelativizeToPublishedPageFolder(suffixedPage)

	var fileInfo os.FileInfo
	if fileInfo, errResponse = nxfsfiles.GetFileInfoIfPathExistOrErrorResponse(publishedPageFullPath); nil != errResponse {
		return errResponse
	}

	if fileInfo.IsDir() {
		return helper.ErrorResponse(http.StatusUnprocessableEntity, "cannot_unpublish_dir", "The received path corresponds to a directory, only pages can be unpublished")
	}

	if implResponse := nxfsfiles.DeleteFile(publishedPageFullPath); implResponse != nil {
		return implResponse
	}

	return nil
}

// addSuffix - receive a string and add the suffix if not present, then return it
func addPageSuffix(value string) (suffixedString string) {
	if strings.HasSuffix(value, pageSuffix) {
		suffixedString = value
	} else {
		suffixedString = value + pageSuffix
	}
	return suffixedString
}
