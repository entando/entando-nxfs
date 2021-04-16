// +build unit

package helper

import (
	"github.com/entando/entando-nxfs/server/model"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

const testBrowsableFsPath = "testdata"
const dirName = "test_dir"
const content = "nice content"

var modelAssembler ModelAssembler
var fileInfo os.FileInfo
var notExistingFileInfo os.FileInfo
var dirPath string

type TestTableDirectory struct {
	paramPath     string
	paramFileInfo os.FileInfo
	expected      *model.DirectoryObject
}

type TestTableFile struct {
	paramPath     string
	paramFileInfo os.FileInfo
	paramContent  string
	expected      *model.FileObject
}

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func setup() {
	_ = os.Mkdir(testBrowsableFsPath, 0755)
	modelAssembler = NewModelAssembler(func() string { return "./" + testBrowsableFsPath })

	dirPath = path.Join(testBrowsableFsPath, dirName)
	_ = os.Mkdir(dirPath, 0755)
	fileInfo, _ = os.Stat(dirPath)
	notExistingFileInfo, _ = os.Stat(dirPath + "/not-existing")
}

func teardown() {
	os.Remove(path.Join(testBrowsableFsPath, dirName))
	os.Remove(testBrowsableFsPath)
}

func TestToDirectoryObject(t *testing.T) {

	testTables := []*TestTableDirectory{
		{testBrowsableFsPath, fileInfo, &model.DirectoryObject{Name: dirName, Path: "."}}, // correct values
		{testBrowsableFsPath, notExistingFileInfo, &model.DirectoryObject{}},              // not existing FileInfo = empty DirectoryObject
		{"not_existing", notExistingFileInfo, &model.DirectoryObject{}},                   // all params not existing = empty DirectoryObject
		{testBrowsableFsPath, nil, &model.DirectoryObject{}},                              // nil FileInfo = empty DirectoryObject
		{"", nil, &model.DirectoryObject{}},                                               // all params nil = empty DirectoryObject
	}

	for _, tTable := range testTables {

		directoryObject := modelAssembler.ToDirectoryObject(tTable.paramPath, tTable.paramFileInfo)

		assert.Equal(t, tTable.expected.Name, directoryObject.Name)
		assert.Equal(t, tTable.expected.Path, directoryObject.Path)
	}
}

func TestToDirectoryObjectFromFilePath(t *testing.T) {

	testTables := []*TestTableDirectory{
		{paramPath: dirPath, expected: &model.DirectoryObject{Name: dirName, Path: "."}}, // correct values
		{paramPath: "not_existing", expected: &model.DirectoryObject{}},                  // all params not existing = empty DirectoryObject
		{paramPath: "", expected: &model.DirectoryObject{}},                              // all params nil = empty DirectoryObject
	}

	for _, tTable := range testTables {

		directoryObject := modelAssembler.ToDirectoryObjectFromFilePath(tTable.paramPath)

		assert.Equal(t, tTable.expected.Name, directoryObject.Name)
		assert.Equal(t, tTable.expected.Path, directoryObject.Path)
	}
}

func TestToFileObject(t *testing.T) {

	testTables := []*TestTableFile{
		{testBrowsableFsPath, fileInfo, content, &model.FileObject{Name: dirName, Path: ".", Content: content}}, // correct values
		{testBrowsableFsPath, fileInfo, "", &model.FileObject{Name: dirName, Path: "."}},                        // no content
		{testBrowsableFsPath, notExistingFileInfo, content, &model.FileObject{}},                                // not existing FileInfo = empty FileObject
		{"not_existing", notExistingFileInfo, content, &model.FileObject{}},                                     // all params not existing = empty FileObject
		{testBrowsableFsPath, nil, content, &model.FileObject{}},                                                // nil FileInfo = empty FileObject
		{"", nil, "", &model.FileObject{}},                                                                      // all params nil = empty FileObject
	}

	for _, tTable := range testTables {

		fileObject := modelAssembler.ToFileObject(tTable.paramPath, tTable.paramFileInfo, tTable.paramContent)

		assert.Equal(t, tTable.expected.Name, fileObject.Name)
		assert.Equal(t, tTable.expected.Path, fileObject.Path)
		assert.Equal(t, tTable.expected.Content, fileObject.Content)
	}
}

func TestToFileObjectFromFilePath(t *testing.T) {

	testTables := []*TestTableFile{
		{paramPath: dirPath, paramContent: content, expected: &model.FileObject{Name: dirName, Path: ".", Content: content}}, // correct values
		{paramPath: dirPath, paramContent: "", expected: &model.FileObject{Name: dirName, Path: "."}},                        // no content
		{paramPath: "not_existing", paramContent: content, expected: &model.FileObject{}},                                    // not existing path = empty FileObject
		{paramPath: "", paramContent: "", expected: &model.FileObject{}},                                                     // all params nil = empty FileObject
	}

	for _, tTable := range testTables {

		fileObject := modelAssembler.ToFileObjectFromFilePath(tTable.paramPath, tTable.paramContent)

		assert.Equal(t, tTable.expected.Name, fileObject.Name)
		assert.Equal(t, tTable.expected.Path, fileObject.Path)
		assert.Equal(t, tTable.expected.Content, fileObject.Content)
	}
}
