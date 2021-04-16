// +build integration

package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	nxsiteman "github.com/entando/entando-nxfs/server"
	"github.com/entando/entando-nxfs/server/controller"
	"github.com/entando/entando-nxfs/server/helper"
	"github.com/entando/entando-nxfs/server/model"
	"github.com/entando/entando-nxfs/server/service"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var defaultApiService controller.DefaultApiServicer
var defaultApiController nxsiteman.Router

type testTable struct {
	paramPath           string
	paramMaxDepth       int
	paramPublishedPages bool
	setupFn             func(string)
	setupFnParam        string
	body                interface{}
	expStatusCode       int
	expRespBodyJsonFile string
	cleanUpFn           func(string)
	cleanUpFnParam      string
}

func createDir(path string) {
	os.Mkdir(filepath.Join(helper.FsBaseDir, path), 0755)
}

func createFile(path string) {
	ioutil.WriteFile(filepath.Join(helper.FsBaseDir, path), nil, 0755)
}

func removeFileOrDir(path string) {
	os.Remove(filepath.Join(helper.FsBaseDir, path))
}

func TestMain(m *testing.M) {
	defaultApiService = service.NewDefaultApiService()
	defaultApiController = controller.NewDefaultApiController(defaultApiService)

	os.Exit(m.Run())
}

func TestBrowseApi(t *testing.T) {

	testTables := []*testTable{
		&testTable{paramPath: ".%252F", paramMaxDepth: 0, paramPublishedPages: false, expStatusCode: http.StatusOK, expRespBodyJsonFile: "browse_api_response.json"},
		&testTable{paramPath: ".%252F", paramMaxDepth: 0, paramPublishedPages: true, expStatusCode: http.StatusOK, expRespBodyJsonFile: "browse_api_pubpages_response.json"},
		&testTable{paramPath: ".%252F", paramMaxDepth: 1, paramPublishedPages: false, expStatusCode: http.StatusOK, expRespBodyJsonFile: "browse_maxdepth_1_api_response.json"},
		&testTable{paramPath: "dir_level_1", paramMaxDepth: 0, paramPublishedPages: false, expStatusCode: http.StatusOK, expRespBodyJsonFile: "browse_api_path_dir_level_1_response.json"},
		&testTable{paramPath: "pages", paramMaxDepth: 0, paramPublishedPages: true, expStatusCode: http.StatusOK, expRespBodyJsonFile: "browse_api_empty_response.json"},
		&testTable{paramPath: "pub_pages", paramMaxDepth: 0, paramPublishedPages: true, setupFn: createDir, setupFnParam: "pub_pages",
			expStatusCode: http.StatusOK, expRespBodyJsonFile: "browse_api_empty_response.json", cleanUpFn: removeFileOrDir, cleanUpFnParam: "pub_pages"},
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/browse/%s?maxdepth=%d&publishedpages=%t", tTable.paramPath, tTable.paramMaxDepth, tTable.paramPublishedPages)

		executeApiRequest(t, tTable, path,
			"/api/nxfs/browse/{EncodedPath}",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsBrowseEncodedPathGet, "GET")
	}
}

func TestBrowseApiWithoutOptionalParamPublishedPages(t *testing.T) {

	testTables := []*testTable{
		&testTable{paramPath: ".%252F", paramMaxDepth: 0, expStatusCode: http.StatusOK, expRespBodyJsonFile: "browse_api_response.json"},
		&testTable{paramPath: ".%252F", paramMaxDepth: 1, expStatusCode: http.StatusOK, expRespBodyJsonFile: "browse_maxdepth_1_api_response.json"},
		&testTable{paramPath: "dir_level_1", paramMaxDepth: 0, expStatusCode: http.StatusOK, expRespBodyJsonFile: "browse_api_path_dir_level_1_response.json"},
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/browse/%s?maxdepth=%d", tTable.paramPath, tTable.paramMaxDepth)

		executeApiRequest(t, tTable, path,
			"/api/nxfs/browse/{EncodedPath}",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsBrowseEncodedPathGet, "GET")
	}
}

func TestGetObjectApi(t *testing.T) {

	testTables := []*testTable{
		&testTable{paramPath: "root_file.txt", expStatusCode: http.StatusOK, expRespBodyJsonFile: "get_object_api_root_file.json"},
		&testTable{paramPath: "not_existing.txt", expStatusCode: http.StatusNotFound},
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/objects/%s", tTable.paramPath)

		executeApiRequest(t, tTable, path,
			"/api/nxfs/objects/{EncodedPath}",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsObjectsEncodedPathGet, "GET")
	}
}

func TestPutObjectApi(t *testing.T) {

	testTables := []*testTable{
		&testTable{paramPath: "new_file.txt", body: &model.FileObject{Type: model.F, Content: "funky soul"}, expStatusCode: http.StatusCreated, expRespBodyJsonFile: "put_object_api_new_file.json"},
		&testTable{paramPath: "new_dir", body: &model.FileObject{Type: model.D}, expStatusCode: http.StatusCreated, expRespBodyJsonFile: "put_object_api_new_dir.json"},
		&testTable{paramPath: "not_existing_dir/not_existing.txt", body: &model.FileObject{Type: model.F, Content: "funky soul"}, expStatusCode: http.StatusNotFound},
		&testTable{paramPath: "new_dir_2", body: &model.FileObject{Type: model.D, Content: "funky soul"}, expStatusCode: http.StatusBadRequest},
		&testTable{paramPath: "new_file_2.txt", body: &model.FileObject{Type: model.F}, expStatusCode: http.StatusBadRequest},
		//TODO add check on page creation
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/objects/%s", tTable.paramPath)

		executeApiRequest(t, tTable, path, "/api/nxfs/objects/{EncodedPath}",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsObjectsEncodedPathPut, "PUT")

		// clean the env
		if tTable.expStatusCode == http.StatusOK || tTable.expStatusCode == http.StatusCreated {
			err := os.Remove(filepath.Join(helper.FsBaseDir, tTable.paramPath))
			if nil != err {
				fmt.Printf(err.Error())
			}
		}
	}
}

func TestDeleteObjectApi(t *testing.T) {

	testTables := []*testTable{
		&testTable{paramPath: "new_file.txt", setupFn: createFile, setupFnParam: "new_file.txt", expStatusCode: http.StatusNoContent},
		&testTable{paramPath: "new_dir", setupFn: createDir, setupFnParam: "new_dir", expStatusCode: http.StatusNoContent},
		&testTable{paramPath: "not_existing.txt", expStatusCode: http.StatusNoContent}, // will works even if the file does not exist
		&testTable{paramPath: "not_existing", expStatusCode: http.StatusNoContent},     // will works even if the file does not exist
		// will return 422 if dir with children
		&testTable{paramPath: "dir_with_children", setupFn: func(path string) {
			os.Mkdir(filepath.Join(helper.FsBaseDir, "dir_with_children"), 0755)
			ioutil.WriteFile(filepath.Join(helper.FsBaseDir, "dir_with_children/child.txt"), nil, 0755)
		}, expStatusCode: http.StatusUnprocessableEntity, cleanUpFn: func(path string) {
			os.Remove(filepath.Join(helper.FsBaseDir, "dir_with_children/child.txt"))
			os.Remove(filepath.Join(helper.FsBaseDir, "dir_with_children"))
		}},
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/objects/%s", tTable.paramPath)

		executeApiRequest(t, tTable, path, "/api/nxfs/objects/{EncodedPath}",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsObjectsEncodedPathDelete, "DELETE")
	}
}

func TestPublishPageApi(t *testing.T) {

	testTables := []*testTable{
		&testTable{paramPath: "not_existing", expStatusCode: http.StatusNotFound},
		&testTable{paramPath: "home.page", expStatusCode: http.StatusOK},
		&testTable{paramPath: "home", expStatusCode: http.StatusOK},
		&testTable{paramPath: "nested", setupFn: createDir, setupFnParam: "pages/nested.page", expStatusCode: http.StatusUnprocessableEntity}, // can't publish dir
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/objects/%s/publish", tTable.paramPath)

		executeApiRequest(t, tTable, path, "/api/nxfs/objects/{EncodedPath}/publish",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsObjectsEncodedPathPublishPost, "POST")
	}
}

func TestUnpublishPageApi(t *testing.T) {

	testTables := []*testTable{
		&testTable{paramPath: "not_existing", expStatusCode: http.StatusNotFound},
		&testTable{paramPath: "home.page", setupFn: createFile, setupFnParam: "pub_pages/home.page", expStatusCode: http.StatusOK, cleanUpFn: createFile, cleanUpFnParam: "pub_pages/home.page"},
		&testTable{paramPath: "home", setupFn: createFile, setupFnParam: "pub_pages/home.page", expStatusCode: http.StatusOK, cleanUpFn: createFile, cleanUpFnParam: "pub_pages/home.page"},
		&testTable{paramPath: "nested", setupFn: createDir, setupFnParam: "pub_pages/nested.page", expStatusCode: http.StatusUnprocessableEntity, cleanUpFn: removeFileOrDir, cleanUpFnParam: "pub_pages/nested.page"}, // can't publish dir
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/objects/%s/unpublish", tTable.paramPath)

		executeApiRequest(t, tTable, path, "/api/nxfs/objects/{EncodedPath}/unpublish",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsObjectsEncodedPathUnpublishPost, "POST")
	}
}

type apiHandlerFn func(w http.ResponseWriter, r *http.Request)

func executeApiRequest(t *testing.T, tTable *testTable, apiUrl string, apiHandleUrl string, apiHandler apiHandlerFn, httpMethod string) {
	if tTable.setupFn != nil {
		tTable.setupFn(tTable.setupFnParam)
	}

	bodyByte, _ := json.Marshal(tTable.body)
	req, err := http.NewRequest(httpMethod, apiUrl, bytes.NewReader(bodyByte))
	if err != nil {
		t.Fatal(err)
	}

	resRecorder := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(apiHandleUrl, http.HandlerFunc(apiHandler))
	router.ServeHTTP(resRecorder, req)

	assertOnStatusCode(t, resRecorder, tTable)

	if tTable.cleanUpFn != nil {
		tTable.cleanUpFn(tTable.cleanUpFnParam)
	}
}

// readResponseFile - read the json file from the testdata folder and return its content as byte array
func readResponseFile(fileName string) []byte {

	content, err := ioutil.ReadFile("testdata/" + fileName)
	if err != nil {
		log.Fatal(err)
	}

	return content
}

// validateResponse - validate the response (error, status code and response body)
func validateResponse(t *testing.T, resRecorder *httptest.ResponseRecorder, tTable *testTable, objType reflect.Type) {

	assertOnStatusCode(t, resRecorder, tTable)

	if tTable.expRespBodyJsonFile != "" {
		actualRes := reflect.New(objType).Elem().Interface()
		expectedRes := reflect.New(objType).Elem().Interface()

		if err := json.Unmarshal(resRecorder.Body.Bytes(), &actualRes); err != nil {
			t.Errorf("Failed unmarshaling response body")
		}
		if err := json.Unmarshal(readResponseFile(tTable.expRespBodyJsonFile), &expectedRes); err != nil {
			t.Errorf("Failed unmarshaling json file")
		}

		switch actualRes.(type) {
		case model.FileObject:
			assertOnFileObjects(t, actualRes.(model.FileObject), expectedRes.(model.FileObject))
		case model.DirectoryObject:
			assertOnDirectoryObjects(t, actualRes.(model.DirectoryObject), expectedRes.(model.DirectoryObject))
		}
	}
}

func assertOnStatusCode(t *testing.T, resRecorder *httptest.ResponseRecorder, tTable *testTable) {
	if resRecorder.Code != tTable.expStatusCode {
		t.Errorf("Expected status code %d, received %d for path %s", tTable.expStatusCode, resRecorder.Code, tTable.paramPath)
	}
}

func assertOnFileObjects(t *testing.T, actual model.FileObject, expected model.FileObject) {

	if actual.Name != expected.Name {
		t.Errorf("Expected FileObject name %s, received %s", actual.Name, expected.Name)
	}
	if actual.Path != expected.Path {
		t.Errorf("Expected FileObject path %s, received %s", actual.Path, expected.Path)
	}
	if actual.Content != expected.Content {
		t.Errorf("Expected FileObject Content %s, received %s", actual.Content, expected.Content)
	}
	if actual.Size != expected.Size {
		t.Errorf("Expected FileObject Size %d, received %d", actual.Size, expected.Size)
	}
	if actual.Type != model.F {
		t.Errorf("Expected FileObject type %s, F expected", actual.Name)
	}
}

func assertOnDirectoryObjects(t *testing.T, actual model.DirectoryObject, expected model.DirectoryObject) {

	if actual.Name != expected.Name {
		t.Errorf("Expected FileObject name %s, received %s", actual.Name, expected.Name)
	}
	if actual.Path != expected.Path {
		t.Errorf("Expected FileObject path %s, received %s", actual.Path, expected.Path)
	}
	if actual.Size != expected.Size {
		t.Errorf("Expected FileObject Size %d, received %d", actual.Size, expected.Size)
	}
	if actual.Type != model.F {
		t.Errorf("Expected FileObject type %s, F expected", actual.Name)
	}
}
