package api

import (
	"encoding/json"
	// "fmt"
	"net/http"

	"github.com/dghubble/sling"
)

type SchemaService struct {
	sling *sling.Sling
}

func newSchemaService(sling *sling.Sling) *SchemaService {
	return &SchemaService{
		sling: sling.Path("schemas"),
	}
}

func (s *SchemaService) Filter(params interface{}) (*json.RawMessage, *http.Response, error) {
	var schemas *json.RawMessage = &json.RawMessage{}
	apiError := new(APIError)
	resp, err := s.sling.New().QueryStruct(params).Receive(schemas, apiError)
	return schemas, resp, relevantError(err, apiError)
}

func (s *SchemaService) List(params *PaginationParams) (*json.RawMessage, *http.Response, error) {
	var schemas *json.RawMessage = &json.RawMessage{}
	apiError := new(APIError)
	resp, err := s.sling.New().QueryStruct(params).Receive(schemas, apiError)
	return schemas, resp, relevantError(err, apiError)
}

// Get schema of a store. Schemas contain settings and data structure of a store
// GET https://api.dostow.com/v1/schemas/get_schema_id
func (s *SchemaService) Get(schemaID string) (*json.RawMessage, *http.Response, error) {
	var schema *json.RawMessage = &json.RawMessage{}
	apiError := new(APIError)
	resp, err := s.sling.New().Get(schemaID).Receive(schema, apiError)
	return schema, resp, relevantError(err, apiError)
}

// Create a new schema.
// POST https://api.dostow.com/v1/schemas
func (o *SchemaService) Create(storeName string, objectBody interface{}) (*json.RawMessage, *http.Response, error) {
	var schema *json.RawMessage = &json.RawMessage{}
	apiError := new(APIError)
	resp, err := o.sling.New().Post("").BodyJSON(objectBody).Receive(schema, apiError)
	return schema, resp, relevantError(err, apiError)
}

// Update the data on an object that already exists.
func (o *SchemaService) Update(schemaID string, objectBody interface{}) (*json.RawMessage, *http.Response, error) {
	var schema *json.RawMessage = &json.RawMessage{}
	apiError := new(APIError)
	resp, err := o.sling.New().Put("schemas/"+schemaID).BodyJSON(objectBody).Receive(schema, apiError)

	return schema, resp, relevantError(err, apiError)
}

// Delete a schema.
// DELETE
func (o *SchemaService) Delete(schemaID string) (*http.Response, error) {

	apiError := new(APIError)
	resp, err := o.sling.New().Delete("schemas/"+schemaID).Receive(nil, apiError)

	return resp, relevantError(err, apiError)
}
