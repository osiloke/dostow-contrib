package api

import (
	"encoding/json"
	// "fmt"
	"net/http"

	"github.com/dghubble/sling"
)

type StoreService struct {
	sling *sling.Sling
}

func newStoreService(sling *sling.Sling) *StoreService {
	return &StoreService{
		sling: sling.Path("store"),
	}
}

func (s *StoreService) List(store string, opts ...Opt) (*json.RawMessage, *http.Response, error) {
	var rows = &json.RawMessage{}
	apiError := new(APIError)
	_s := s.sling.New()
	for _, opt := range opts {
		_s = opt(_s)
	}
	resp, err := _s.Path("store/"+store).Receive(rows, apiError)
	return rows, resp, relevantError(err, apiError)
}
func (s *StoreService) Get(store, id string, data interface{}) error {
	apiError := new(APIError)
	_s := s.sling.New().Get("store/" + store + "/" + id)
	_, err := _s.Receive(data, apiError)
	return relevantError(err, apiError)
}
func (s *StoreService) GetRaw(store, id string, opts ...Opt) (*json.RawMessage, error) {
	var result = &json.RawMessage{}
	apiError := new(APIError)
	_s := s.sling.New()
	for _, opt := range opts {
		_s = opt(_s)
	}

	resp, err := _s.Get("store/"+store+"/"+id).Receive(result, apiError)
	if resp.StatusCode == 404 {
		apiError.Status = "404"
		apiError.Message = "not found"
		return nil, apiError
	}
	return result, relevantError(err, apiError)
}
func (s *StoreService) Create(store string, data interface{}, opts ...Opt) (*json.RawMessage, error) {
	var result = &json.RawMessage{}
	apiError := new(APIError)
	_s := s.sling.New()
	for _, opt := range opts {
		_s = opt(_s)
	}
	_, err := _s.Post("store/"+store).BodyJSON(data).Receive(result, apiError)
	return result, relevantError(err, apiError)
}
func (s *StoreService) BulkCreate(store string, data interface{}, opts ...Opt) (*json.RawMessage, error) {
	var result = &json.RawMessage{}
	apiError := new(APIError)
	_s := s.sling.New()
	for _, opt := range opts {
		_s = opt(_s)
	}
	_, err := _s.Post("store/"+store+"/bulk").BodyJSON(data).Receive(result, apiError)
	return result, relevantError(err, apiError)
}
func (s *StoreService) Update(store, id string, data interface{}, opts ...Opt) (*json.RawMessage, error) {
	var result = &json.RawMessage{}
	apiError := new(APIError)
	_s := s.sling.New()
	for _, opt := range opts {
		_s = opt(_s)
	}
	_, err := _s.Put("store/"+store+"/"+id).BodyJSON(data).Receive(result, apiError)
	return result, relevantError(err, apiError)
}

func (s *StoreService) Remove(store, id string) (*json.RawMessage, error) {
	var result = &json.RawMessage{}
	apiError := new(APIError)
	_s := s.sling.New().Delete("store/" + store + "/" + id)
	_, err := _s.Receive(result, apiError)
	return result, relevantError(err, apiError)
}

//ListResult contains a query result
type ListResult struct {
	Total int64 `json:"total_count"`
}

//Size gets the size of store
func (s *StoreService) Size(store string) (cnt int64, err error) {
	res := new(ListResult)
	raw, rsp, err := s.List(store, s.Query(PaginationParams{Size: 1}))
	if err == nil {
		if err = json.Unmarshal(*raw, res); err == nil {
			cnt = res.Total
		}

	} else if rsp.StatusCode == 404 {
		if err = json.Unmarshal(*raw, res); err == nil {
			cnt = res.Total
		}
	}
	return
}

//Clear a store
func (s *StoreService) Clear(store string) (*json.RawMessage, *http.Response, error) {
	var result = &json.RawMessage{}
	apiError := new(APIError)
	_s := s.sling.New().Delete("clear_store/" + store)
	rsp, err := _s.Receive(result, apiError)
	return result, rsp, relevantError(err, apiError)
}

func (s *StoreService) Authorize(token string) func(sl *sling.Sling) *sling.Sling {
	return Authorize(token)
}
func (s *StoreService) Query(q interface{}) func(sl *sling.Sling) *sling.Sling {
	return Query(q)
}
