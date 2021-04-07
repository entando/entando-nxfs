package nxsiteman

import (
	"net/http"
	"os"
	"strings"
)

const publishedPageSuffix = ".page"
const draftPageSuffix = publishedPageSuffix + ".draft"

// publishPage - publish the received draft page and return an error ImplResponse if an error occurs, nil otherwise
func publishPage(encodedDraftPagePath string) (errorResp *ImplResponse) {

	var pageFileInfo os.FileInfo

	// decode path
	decodedPath, errResponse := decodePath(encodedDraftPagePath)
	if errResponse != nil {
		return errResponse
	}

	suffixedDraftPage := manageDraftPageSuffix(decodedPath)
	suffixedPublishedPage := managePublishPageSuffix(decodedPath)

	// check if file exist as draft in the correct folder or error
	pageFileInfo, errResponse = getDraftPageInfoIfExistOrErrorResponse(suffixedDraftPage)
	if errResponse != nil {
		return errResponse
	}

	// if dir error
	if pageFileInfo.IsDir() {
		return ErrorResponse(http.StatusUnprocessableEntity, "cannot_publish_dir", "The received path corresponds to a directory, only pages can be published")
	}

	draftPageFullPath := relativizeToDraftPageFolder(suffixedDraftPage)
	publishedPageFullPath := relativizeToPublishedPageFolder(suffixedPublishedPage)
	return copyFileTo(draftPageFullPath, publishedPageFullPath)
}

// manageDraftPageSuffix - receive a draft page file path and add the page.draft suffix if needed, then return it
func manageDraftPageSuffix(draftPagePath string) string {
	return addSuffix(draftPagePath, draftPageSuffix)
}

// managePublishPageSuffix - receive a published page file path and ensure that it contains the right suffix, then return it
func managePublishPageSuffix(publishedPagePath string) (suffixedString string) {
	suffixedString = removeSuffix(publishedPagePath, draftPageSuffix)
	return addSuffix(suffixedString, publishedPageSuffix)
}

// addSuffix - receive a string and add the suffix if not present, then return it
func addSuffix(value string, suffix string) (suffixedString string) {
	if strings.HasSuffix(value, suffix) {
		suffixedString = value
	} else {
		suffixedString = value + suffix
	}
	return suffixedString
}

// removeSuffix - receive a string and remove the suffix if present, then return it
func removeSuffix(value string, suffix string) (suffixedString string) {
	return strings.TrimSuffix(value, suffix)
}
