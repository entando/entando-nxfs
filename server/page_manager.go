package nxsiteman

import (
	"net/http"
	"os"
	"strings"
)

const pageSuffix = ".page"

// publishPage - publish the received draft page and return an error ImplResponse if an error occurs, nil otherwise
func publishPage(encodedDraftPagePath string) (errorResp *ImplResponse) {

	var pageFileInfo os.FileInfo

	// decode path
	decodedPath, errResponse := decodePath(encodedDraftPagePath)
	if errResponse != nil {
		return errResponse
	}

	suffixedPage := addPageSuffix(decodedPath)

	// check if file exist as draft in the correct folder or error
	pageFileInfo, errResponse = getDraftPageInfoIfExistOrErrorResponse(suffixedPage)
	if errResponse != nil {
		return errResponse
	}

	// if dir error
	if pageFileInfo.IsDir() {
		return ErrorResponse(http.StatusUnprocessableEntity, "cannot_publish_dir", "The received path corresponds to a directory, only pages can be published")
	}

	draftPageFullPath := relativizeToDraftPageFolder(suffixedPage)
	publishedPageFullPath := relativizeToPublishedPageFolder(suffixedPage)
	return copyFileTo(draftPageFullPath, publishedPageFullPath)
}

// unpublishPage - unpublish the received published page and return an error ImplResponse if an error occurs, nil otherwise
func unpublishPage(encodedPublishedPagePath string) (errorResp *ImplResponse) {

	// decode path
	decodedPath, errResponse := decodePath(encodedPublishedPagePath)
	if errResponse != nil {
		return errResponse
	}

	suffixedPage := addPageSuffix(decodedPath)
	publishedPageFullPath := relativizeToPublishedPageFolder(suffixedPage)

	if _, implResponse := getFileInfoIfPathExistOrErrorResponse(publishedPageFullPath); nil != implResponse {
		return implResponse
	}

	if implResponse := deleteFile(publishedPageFullPath); implResponse != nil {
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
