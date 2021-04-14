/*
 * NxFs
 *
 * Simple file access APIs for the Entando Nx subsystem
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package controller

import (
	"encoding/json"
	nxsiteman "github.com/entando/entando-nxfs/server"
	"github.com/entando/entando-nxfs/server/model"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// A DefaultApiController binds http requests to an api service and writes the service results to the http response
type DefaultApiController struct {
	service DefaultApiServicer
}

// NewDefaultApiController creates a default api controller
func NewDefaultApiController(s DefaultApiServicer) nxsiteman.Router {
	return &DefaultApiController{service: s}
}

// Routes returns all of the api route for the DefaultApiController
func (c *DefaultApiController) Routes() nxsiteman.Routes {
	return nxsiteman.Routes{
		{
			"ApiNxfsBrowseEncodedPathGet",
			strings.ToUpper("Get"),
			"/api/nxfs/browse/{EncodedPath}",
			c.ApiNxfsBrowseEncodedPathGet,
		},
		{
			"ApiNxfsObjectsEncodedPathDelete",
			strings.ToUpper("Delete"),
			"/api/nxfs/objects/{EncodedPath}",
			c.ApiNxfsObjectsEncodedPathDelete,
		},
		{
			"ApiNxfsObjectsEncodedPathGet",
			strings.ToUpper("Get"),
			"/api/nxfs/objects/{EncodedPath}",
			c.ApiNxfsObjectsEncodedPathGet,
		},
		{
			"ApiNxfsObjectsEncodedPathPublishPost",
			strings.ToUpper("Post"),
			"/api/nxfs/objects/{EncodedPath}/publish",
			c.ApiNxfsObjectsEncodedPathPublishPost,
		},
		{
			"ApiNxfsObjectsEncodedPathPut",
			strings.ToUpper("Put"),
			"/api/nxfs/objects/{EncodedPath}",
			c.ApiNxfsObjectsEncodedPathPut,
		},
		{
			"ApiNxfsObjectsEncodedPathUnpublishPost",
			strings.ToUpper("Post"),
			"/api/nxfs/objects/{EncodedPath}/unpublish",
			c.ApiNxfsObjectsEncodedPathUnpublishPost,
		},
	}
}

// ApiNxfsBrowseEncodedPathGet - Gets the list of objects in a directory
func (c *DefaultApiController) ApiNxfsBrowseEncodedPathGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	encodedPath := params["EncodedPath"]
	maxdepth, err := nxsiteman.ParseInt32Parameter(query.Get("maxdepth"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	publishedPages, err := nxsiteman.ParseBoolParameter(query.Get("publishedpages"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := c.service.ApiNxfsBrowseEncodedPathGet(r.Context(), encodedPath, maxdepth, publishedPages)
	//If an error occured, encode the error with the status code
	if err != nil {
		nxsiteman.EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	nxsiteman.EncodeJSONResponse(result.Body, &result.Code, w)

}

// ApiNxfsObjectsEncodedPathDelete - Deletes an object
func (c *DefaultApiController) ApiNxfsObjectsEncodedPathDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	encodedPath := params["EncodedPath"]
	result, err := c.service.ApiNxfsObjectsEncodedPathDelete(r.Context(), encodedPath)
	//If an error occured, encode the error with the status code
	if err != nil {
		nxsiteman.EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	nxsiteman.EncodeJSONResponse(result.Body, &result.Code, w)

}

// ApiNxfsObjectsEncodedPathGet - Gets an object
func (c *DefaultApiController) ApiNxfsObjectsEncodedPathGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	encodedPath := params["EncodedPath"]
	result, err := c.service.ApiNxfsObjectsEncodedPathGet(r.Context(), encodedPath)
	//If an error occured, encode the error with the status code
	if err != nil {
		nxsiteman.EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	nxsiteman.EncodeJSONResponse(result.Body, &result.Code, w)

}

// ApiNxfsObjectsEncodedPathPublishPost - Publishes a page
func (c *DefaultApiController) ApiNxfsObjectsEncodedPathPublishPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	encodedPath := params["EncodedPath"]
	result, err := c.service.ApiNxfsObjectsEncodedPathPublishPost(r.Context(), encodedPath)
	//If an error occured, encode the error with the status code
	if err != nil {
		nxsiteman.EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	nxsiteman.EncodeJSONResponse(result.Body, &result.Code, w)

}

// ApiNxfsObjectsEncodedPathPut - Creates or updates an object
func (c *DefaultApiController) ApiNxfsObjectsEncodedPathPut(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	encodedPath := params["EncodedPath"]
	fileObject := &model.FileObject{}
	if err := json.NewDecoder(r.Body).Decode(&fileObject); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := c.service.ApiNxfsObjectsEncodedPathPut(r.Context(), encodedPath, *fileObject)
	//If an error occured, encode the error with the status code
	if err != nil {
		nxsiteman.EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	nxsiteman.EncodeJSONResponse(result.Body, &result.Code, w)

}

// ApiNxfsObjectsEncodedPathUnpublishPost - Unpublishes a page
func (c *DefaultApiController) ApiNxfsObjectsEncodedPathUnpublishPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	encodedPath := params["EncodedPath"]
	result, err := c.service.ApiNxfsObjectsEncodedPathUnpublishPost(r.Context(), encodedPath)
	//If an error occured, encode the error with the status code
	if err != nil {
		nxsiteman.EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	nxsiteman.EncodeJSONResponse(result.Body, &result.Code, w)

}
