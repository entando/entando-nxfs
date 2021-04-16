// +build unit

package helper

import (
	"github.com/entando/entando-nxfs/server/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSuccessResponse(t *testing.T) {

	type TestTableSuccessResponse struct {
		Code int
		Body interface{}
	}

	type TestResult struct {
		My      string
		Structz string
	}

	testTables := []*TestTableSuccessResponse{
		{Code: 10, Body: "result"},
		{Code: -1, Body: "result"},
		{Code: 10, Body: &TestResult{My: "Nice", Structz: "Feature"}},
		{Code: 10},
		{Body: "result"},
		{}, // TODO what should do? return error?
	}

	for _, tTable := range testTables {
		response := SuccessResponse(tTable.Code, tTable.Body)
		assert.Equal(t, tTable.Code, response.Code)
		assert.Equal(t, tTable.Body, response.Body)
	}
}

func TestErrorResponse(t *testing.T) {

	type TestTableErrorResponse struct {
		Code         int
		ErrorCode    string
		ErrorMessage string
	}

	testTables := []*TestTableErrorResponse{
		{Code: 10, ErrorCode: "code", ErrorMessage: "mex"},
		{Code: -1, ErrorCode: "code"},
		{Code: 10, ErrorMessage: "mex"},
		{Code: 10},
		{ErrorCode: "code", ErrorMessage: "mex"},
		{ErrorCode: "code"},
		{ErrorMessage: "mex"},
		{}, // TODO what should do? return error?
	}

	for _, tTable := range testTables {
		response := ErrorResponse(tTable.Code, tTable.ErrorCode, tTable.ErrorMessage)
		result := response.Body.(*model.Result)
		assert.Equal(t, tTable.Code, response.Code)
		assert.Equal(t, tTable.ErrorCode, result.Code)
		assert.Equal(t, tTable.ErrorMessage, result.Message)
	}
}

func TestGetBrowsableFsRootPath(t *testing.T) {

	// TODO add test for env var when it will be implemented

	actual := GetBrowsableFsRootPath()
	expected := "./browsableFS"

	assert.Equal(t, actual, expected)
}

func TestGetPublishedPagesPath(t *testing.T) {

	// TODO add test for env var when it will be implemented

	actual := GetPublishedPagesPath()
	expected := "browsableFS/pub_pages"

	assert.Equal(t, actual, expected)
}

func TestGetDraftPagesPath(t *testing.T) {

	// TODO add test for env var when it will be implemented

	actual := GetDraftPagesPath()
	expected := "browsableFS/pages"

	assert.Equal(t, actual, expected)
}
