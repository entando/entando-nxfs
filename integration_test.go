// +build integration

package main_test

import (
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

type TestTable struct {
	ParamPath           string
	ParamMaxDepth       int
	ParamPublishedPages bool
	SetupFn             func(string)
	SetupFnParam        string
	Body                interface{}
	ExpStatusCode       int
	ExpRespBodyJsonFile string
	CleanUpFn           func(string)
	CleanUpFnParam      string
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

	testTables := []*TestTable{
		&TestTable{ParamPath: ".%252F", ParamMaxDepth: 0, ParamPublishedPages: false, ExpStatusCode: http.StatusOK, ExpRespBodyJsonFile: "browse_api_response.json"},
		&TestTable{ParamPath: ".%252F", ParamMaxDepth: 0, ParamPublishedPages: true, ExpStatusCode: http.StatusOK, ExpRespBodyJsonFile: "browse_api_pubpages_response.json"},
		&TestTable{ParamPath: ".%252F", ParamMaxDepth: 1, ParamPublishedPages: false, ExpStatusCode: http.StatusOK, ExpRespBodyJsonFile: "browse_maxdepth_1_api_response.json"},
		&TestTable{ParamPath: "dir_level_1", ParamMaxDepth: 0, ParamPublishedPages: false, ExpStatusCode: http.StatusOK, ExpRespBodyJsonFile: "browse_api_path_dir_level_1_response.json"},
		&TestTable{ParamPath: "pages", ParamMaxDepth: 0, ParamPublishedPages: true, ExpStatusCode: http.StatusOK, ExpRespBodyJsonFile: "browse_api_empty_response.json"},
		&TestTable{ParamPath: "pub_pages", ParamMaxDepth: 0, ParamPublishedPages: true, ExpStatusCode: http.StatusOK, ExpRespBodyJsonFile: "browse_api_empty_response.json"},
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/browse/%s?maxdepth=%d&publishedpages=%t", tTable.ParamPath, tTable.ParamMaxDepth, tTable.ParamPublishedPages)

		executeApiRequest(t, tTable, path,
			"/api/nxfs/browse/{EncodedPath}",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsBrowseEncodedPathGet, "GET")
	}
}

func TestGetObjectApi(t *testing.T) {

	testTables := []*TestTable{
		&TestTable{ParamPath: "root_file.txt", ExpStatusCode: http.StatusOK, ExpRespBodyJsonFile: "get_object_api_root_file.json"},
		&TestTable{ParamPath: "not_existing.txt", ExpStatusCode: http.StatusNotFound},
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/objects/%s", tTable.ParamPath)

		executeApiRequest(t, tTable, path,
			"/api/nxfs/objects/{EncodedPath}",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsObjectsEncodedPathGet, "GET")
	}
}

func TestPutObjectApi(t *testing.T) {

	testTables := []*TestTable{
		&TestTable{ParamPath: "new_file.txt", Body: &model.FileObject{Type: model.F, Content: "funky soul"}, ExpStatusCode: http.StatusCreated, ExpRespBodyJsonFile: "put_object_api_new_file.json"},
		&TestTable{ParamPath: "new_dir", Body: &model.FileObject{Type: model.D}, ExpStatusCode: http.StatusCreated, ExpRespBodyJsonFile: "put_object_api_new_dir.json"},
		&TestTable{ParamPath: "not_existing_dir/not_existing.txt", Body: &model.FileObject{Type: model.F, Content: "funky soul"}, ExpStatusCode: http.StatusNotFound},
		&TestTable{ParamPath: "new_dir_2", Body: &model.FileObject{Type: model.D, Content: "funky soul"}, ExpStatusCode: http.StatusBadRequest},
		&TestTable{ParamPath: "new_file_2.txt", Body: &model.FileObject{Type: model.F}, ExpStatusCode: http.StatusBadRequest},
		//TODO add check on page creation
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/objects/%s", tTable.ParamPath)

		executeApiRequest(t, tTable, path, "/api/nxfs/objects/{EncodedPath}",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsObjectsEncodedPathPut, "PUT")

		// clean the env
		if tTable.ExpStatusCode == http.StatusOK || tTable.ExpStatusCode == http.StatusCreated {
			err := os.Remove(filepath.Join(helper.FsBaseDir, tTable.ParamPath))
			if nil != err {
				fmt.Printf(err.Error())
			}
		}
	}
}

func TestDeleteObjectApi(t *testing.T) {

	testTables := []*TestTable{
		&TestTable{ParamPath: "new_file.txt", SetupFn: createFile, SetupFnParam: "new_file.txt", ExpStatusCode: http.StatusNoContent},
		&TestTable{ParamPath: "new_dir", SetupFn: createDir, SetupFnParam: "new_dir", ExpStatusCode: http.StatusNoContent},
		&TestTable{ParamPath: "not_existing.txt", ExpStatusCode: http.StatusNoContent}, // will works even if the file does not exist
		&TestTable{ParamPath: "not_existing", ExpStatusCode: http.StatusNoContent},     // will works even if the file does not exist
		// will return 422 if dir with children
		&TestTable{ParamPath: "dir_with_children", SetupFn: func(path string) {
			os.Mkdir(filepath.Join(helper.FsBaseDir, "dir_with_children"), 0755)
			ioutil.WriteFile(filepath.Join(helper.FsBaseDir, "dir_with_children/child.txt"), nil, 0755)
		}, ExpStatusCode: http.StatusUnprocessableEntity, CleanUpFn: func(path string) {
			os.Remove(filepath.Join(helper.FsBaseDir, "dir_with_children/child.txt"))
			os.Remove(filepath.Join(helper.FsBaseDir, "dir_with_children"))
		}},
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/objects/%s", tTable.ParamPath)

		executeApiRequest(t, tTable, path, "/api/nxfs/objects/{EncodedPath}",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsObjectsEncodedPathDelete, "DELETE")
	}
}

func TestPublishPageApi(t *testing.T) {

	testTables := []*TestTable{
		&TestTable{ParamPath: "not_existing", ExpStatusCode: http.StatusNotFound},
		&TestTable{ParamPath: "home.page", ExpStatusCode: http.StatusOK},
		&TestTable{ParamPath: "home", ExpStatusCode: http.StatusOK},
		&TestTable{ParamPath: "nested", SetupFn: createDir, SetupFnParam: "pages/nested.page", ExpStatusCode: http.StatusUnprocessableEntity}, // can't publish dir
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/objects/%s/publish", tTable.ParamPath)

		executeApiRequest(t, tTable, path, "/api/nxfs/objects/{EncodedPath}/publish",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsObjectsEncodedPathPublishPost, "POST")
	}
}

func TestUnpublishPageApi(t *testing.T) {

	testTables := []*TestTable{
		&TestTable{ParamPath: "not_existing", ExpStatusCode: http.StatusNotFound},
		&TestTable{ParamPath: "home.page", SetupFn: createFile, SetupFnParam: "pub_pages/home.page", ExpStatusCode: http.StatusOK, CleanUpFn: createFile, CleanUpFnParam: "pub_pages/home.page"},
		&TestTable{ParamPath: "home", SetupFn: createFile, SetupFnParam: "pub_pages/home.page", ExpStatusCode: http.StatusOK, CleanUpFn: createFile, CleanUpFnParam: "pub_pages/home.page"},
		&TestTable{ParamPath: "nested", SetupFn: createDir, SetupFnParam: "pub_pages/nested.page", ExpStatusCode: http.StatusUnprocessableEntity, CleanUpFn: removeFileOrDir, CleanUpFnParam: "pub_pages/nested.page"}, // can't publish dir
	}

	for _, tTable := range testTables {

		path := fmt.Sprintf("/api/nxfs/objects/%s/unpublish", tTable.ParamPath)

		executeApiRequest(t, tTable, path, "/api/nxfs/objects/{EncodedPath}/unpublish",
			defaultApiController.(*controller.DefaultApiController).ApiNxfsObjectsEncodedPathUnpublishPost, "POST")
	}
}

type apiHandlerFn func(w http.ResponseWriter, r *http.Request)

func executeApiRequest(t *testing.T, tTable *TestTable, apiUrl string, apiHandleUrl string, apiHandler apiHandlerFn, httpMethod string) {
	if tTable.SetupFn != nil {
		tTable.SetupFn(tTable.SetupFnParam)
	}

	req, err := http.NewRequest(httpMethod, apiUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	resRecorder := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(apiHandleUrl, http.HandlerFunc(apiHandler))
	router.ServeHTTP(resRecorder, req)

	assertOnStatusCode(t, resRecorder, tTable)

	if tTable.CleanUpFn != nil {
		tTable.CleanUpFn(tTable.CleanUpFnParam)
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
func validateResponse(t *testing.T, resRecorder *httptest.ResponseRecorder, tTable *TestTable, objType reflect.Type) {

	assertOnStatusCode(t, resRecorder, tTable)

	if tTable.ExpRespBodyJsonFile != "" {
		actualRes := reflect.New(objType).Elem().Interface()
		expectedRes := reflect.New(objType).Elem().Interface()

		if err := json.Unmarshal(resRecorder.Body.Bytes(), &actualRes); err != nil {
			t.Errorf("Failed unmarshaling response body")
		}
		if err := json.Unmarshal(readResponseFile(tTable.ExpRespBodyJsonFile), &expectedRes); err != nil {
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

func assertOnStatusCode(t *testing.T, resRecorder *httptest.ResponseRecorder, tTable *TestTable) {
	if resRecorder.Code != tTable.ExpStatusCode {
		t.Errorf("Expected status code %d, received %d for path %s", tTable.ExpStatusCode, resRecorder.Code, tTable.ParamPath)
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
