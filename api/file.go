package api

import (
	"encoding/json"
	// "fmt"
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/dghubble/sling"
)

type FileService struct {
	sling *sling.Sling
}

func newFileService(sling *sling.Sling) *FileService {
	return &FileService{
		sling: sling.Path("upload"),
	}
}

// http://stackoverflow.com/questions/20205796/golang-post-data-using-the-content-type-multipart-form-data
func (s *FileService) Create(store, key, field, filename string, file io.Reader, opts ...Opt) (*http.Response, *json.RawMessage, error) {
	var result *json.RawMessage = &json.RawMessage{}
	apiError := new(APIError)
	_s := s.sling.New()
	for _, opt := range opts {
		_s = opt(_s)
	}
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile(field, filename)
	if err != nil {
		return nil, nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return nil, nil, err
	}
	w.Close()
	ct := w.FormDataContentType()
	rsp, err := _s.Post("upload/"+store+"/"+key+"/"+field).Set(
		"Content-Type",
		ct).Body(&b).Receive(
		result,
		apiError)

	return rsp, result, relevantError(err, apiError)
}
func (s *FileService) Authorize(token string) func(sl *sling.Sling) *sling.Sling {
	return Authorize(token)
}
func (s *FileService) Query(q interface{}) func(sl *sling.Sling) *sling.Sling {
	return Query(q)
}
