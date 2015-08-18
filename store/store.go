package store

import (
	"errors"
	"github.com/franela/goreq"
	"github.com/mgutz/logxi/v1"
	. "github.com/osiloke/gostore"
	"net/url"
)

var logger = log.New("gostore.dostow")

type ServerError struct {
	code int
	msg  string
}

func (e ServerError) Error() string {
	return e.msg
}

func newServerError(code int, msg string) ServerError {
	return ServerError{code, msg}
}

type DostowRows struct {
}

func handleError(err error) error {
	if e, ok := err.(*url.Error); ok {
		//e.Err == *net.OpError
		if serr, ok := err.(*goreq.Error); ok {
			if serr.Timeout() {
				return newServerError(500, "Server timeout")
			}
			return newServerError(500, err.Error())
		}
		if e.Err.Error() == "no such host" {
			return newServerError(500, "Server Unreachable")
		}
		return newServerError(500, "Server Unreachable")
	}
	return newServerError(500, "Unable to complete action")
}

func handleResponse(res *goreq.Response, dst interface{}) error {
	switch res.StatusCode {
	case 200:
		res.Body.FromJsonTo(&dst)
		return nil
	case 500, 400, 401:
		return newServerError(res.StatusCode, "Unable to perform action due to server error")
	default:
	}
	return errors.New("Cannot perfrom action")
}

func (s DostowRows) Next(dst interface{}) (bool, error) {

	return false, nil
}

func (s DostowRows) Close() {
}

type Dostow struct {
	url   string
	key   string
	debug bool
}

func (s Dostow) CreateDatabase() (err error) {
	return errors.New("not implemented")
}

func (s Dostow) GetStore() interface{} {
	return map[string]string{"url": s.url, "key": s.key}
}

func (s Dostow) CreateTable(store string, sample interface{}) (err error) {
	return nil
}
func (s Dostow) req(path string) *goreq.Request {
	req := goreq.Request{
		Uri:         s.url + "/store/" + path,
		Accept:      "application/json",
		ContentType: "application/json",
		Compression: goreq.Gzip(),
		ShowDebug:   s.debug,
	}
	req.AddHeader("X-DOSTOW-GROUP-ACCESS-KEY", s.key)
	return &req
}
func (s Dostow) get(path string) *goreq.Request {
	return s.req(path)
}

func (s Dostow) post(path string) *goreq.Request {
	req := s.req(path)
	req.Method = "POST"
	return req
}

func (s Dostow) delete(path string) *goreq.Request {
	req := s.req(path)
	req.Method = "DELETE"
	return req
}

func (s Dostow) put(path string) *goreq.Request {
	req := s.req(path)
	req.Method = "PUT"
	return req
}

func (s Dostow) All(count int, skip int, store string) (rrows ObjectRows, err error) {
	return nil, errors.New("not implemented")
}

func (s Dostow) AllCursor(store string) (ObjectRows, error) {
	return nil, errors.New("not implemented")
}

//Before will retrieve all old rows that were created before the row with id was created
// [1, 2, 3, 4], before 2 will return [3, 4]
//r.db('worksmart').table('store').orderBy({'index': r.desc('id')}).filter(r.row('schemas')
// .eq('osiloke_tsaboin_silverbird').and(r.row('id').lt('55b54e93f112a16514000057')))
// .pluck('schemas', 'id','tid', 'timestamp', 'created_at').limit(100)
func (s Dostow) Before(id string, count int, skip int, store string) (rows ObjectRows, err error) {
	return nil, errors.New("not implemented")
}

func (s Dostow) FilterBefore(id string, filter map[string]interface{}, count int, skip int, store string) (rows ObjectRows, err error) {
	return nil, errors.New("not implemented")
}

func (s Dostow) FilterBeforeCount(id string, filter map[string]interface{}, count int, skip int, store string) (int64, error) {
	return 0, errors.New("not implemented")
}

func (s Dostow) FilterSince(id string, filter map[string]interface{}, count int, skip int, store string) (rows ObjectRows, err error) {
	return nil, errors.New("not implemented")
}

//This will retrieve all new rows that were created since the row with id was created
// [1, 2, 3, 4], since 2 will return [1]
func (s Dostow) Since(id string, count, skip int, store string) (rrows ObjectRows, err error) {
	return nil, errors.New("not implemented")
}

func (s Dostow) Get(id, store string, dst interface{}) (err error) {
	return errors.New("not implemented")
}

func (s Dostow) Save(store string, src interface{}) (key string, err error) {
	var msg map[string]interface{}
	req := s.post(store)
	req.Body = src
	res, err := req.Do()
	if err != nil {
		log.Error("Error:", "data", src, "err", err)
		return "", handleError(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	switch res.StatusCode {
	case 200:
		res.Body.FromJsonTo(&msg)
		return msg["id"].(string), nil
	case 404:
		body, _ := res.Body.ToString()
		return "", newServerError(res.StatusCode, body)
	case 500, 400, 401:
		body, _ := res.Body.ToString()
		return "", newServerError(res.StatusCode, body)
	default:
		return "", errors.New("Cannot perfrom action")
	}
}

func (s Dostow) Update(id string, store string, src interface{}) (err error) {
	return errors.New("not implemented")
}

func (s Dostow) Delete(id string, store string) (err error) {
	return errors.New("not implemented")
}

func (s Dostow) Stats(store string) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (s Dostow) GetByField(name, val, store string, dst interface{}) (err error) {
	return errors.New("not implemented")
}

func (s Dostow) FilterGet(filter map[string]interface{}, store string, dst interface{}) (err error) {
	req := s.get(store)
	req.QueryString = filter
	res, err := req.Do()
	if err != nil {
		return handleError(err)
	}
	return handleResponse(res, dst)
}

func (s Dostow) FilterGetAll(filter map[string]interface{}, count int, skip int, store string) (rrows ObjectRows, err error) {
	return nil, errors.New("not implemented")
}

func (s Dostow) FilterDelete(filter map[string]interface{}, store string) (err error) {
	return errors.New("not implemented")
}

func (s Dostow) FilterCount(filter map[string]interface{}, store string) (int64, error) {
	return 0, errors.New("not implemented")
}

func (s Dostow) GetByFieldsByField(name, val, store string, fields []string, dst interface{}) (err error) {
	return errors.New("not implemented")
}

func (s Dostow) Close() {

}

func NewStore(apiUrl, accessKey string) Dostow {
	s := Dostow{url: apiUrl, key: accessKey, debug: true}
	return s
}
