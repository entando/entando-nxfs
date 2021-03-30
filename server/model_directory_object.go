/*
 * NxFs
 *
 * Simple file access APIs for the Entando Nx subsystem
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package nxsiteman

type DirectoryObject struct {
	Id int64 `json:"id"`

	Path string `json:"path,omitempty"`

	Size int32 `json:"size,omitempty"`

	Type ObjectType `json:"type,omitempty"`

	Created ActionLog `json:"_created,omitempty"`

	Updated ActionLog `json:"_updated,omitempty"`
}
