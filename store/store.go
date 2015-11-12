package store

import (
	// "bytes"
	"encoding/json"
	"errors"
	"github.com/beefsack/go-rate"
	"github.com/gosexy/to"
	"github.com/mgutz/logxi/v1"
	"github.com/mitchellh/mapstructure"
	. "github.com/osiloke/gostore"
	"strings"
	// "github.com/smallnest/goreq"
	"github.com/ddliu/go-httpclient"
	//https://github.com/sethgrid/pester //use pester with go-httpclient
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
	// "os"
	"path/filepath"
)

// TODO: add github.com/cenkalti/backoff
const (
	USERAGENT       = "Dostow Store API"
	TIMEOUT         = 30
	CONNECT_TIMEOUT = 5
)

var rl = rate.New(1, time.Second)
var logger = log.New("gostore.dostow")

type ServerError struct {
	code int
	msg  string
}

func (s ServerError) Code() int {
	return s.code
}
func (s ServerError) Msg() string {
	return s.msg
}

type DataList struct {
	TotalCount int           `json:"total_count"`
	Count      int           `json:"count"`
	Data       []interface{} `json:"data"`
}

func (e ServerError) Error() string {
	return e.msg
}

func newServerError(code int, msg string) ServerError {
	return ServerError{code, msg}
}

type DostowRows struct {
	data []interface{}
}

func (s DostowRows) Next(dst interface{}) (bool, error) {

	return false, nil
}

func (s DostowRows) Close() {
}

type Dostow struct {
	url    string
	key    string
	client *http.Client
	debug  bool
}

func (s Dostow) Options() map[string]interface{} {
	return map[string]interface{}{
		"url": s.url,
	}
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

func (s Dostow) FilterBefore(id string, filter map[string]interface{}, count int, skip int, store string, opts ObjectStoreOptions) (rows ObjectRows, err error) {
	return nil, errors.New("not implemented")
}

func (s Dostow) FilterBeforeCount(id string, filter map[string]interface{}, count int, skip int, store string, opts ObjectStoreOptions) (int64, error) {
	return 0, errors.New("not implemented")
}

func (s Dostow) FilterSince(id string, filter map[string]interface{}, count int, skip int, store string, opts ObjectStoreOptions) (rows ObjectRows, err error) {
	return nil, errors.New("not implemented")
}

//This will retrieve all new rows that were created since the row with id was created
// [1, 2, 3, 4], since 2 will return [1]
func (s Dostow) Since(id string, count, skip int, store string) (rrows ObjectRows, err error) {
	return nil, errors.New("not implemented")
}

func (s *Dostow) post(url string, data string) (*httpclient.Response, error) {
	// Post with files should be sent as multipart.
	headers := make(map[string]string)
	headers["X-DOSTOW-GROUP-ACCESS-KEY"] = s.key
	headers["Content-Type"] = "application/json"
	body := strings.NewReader(data)
	// httpclient.NewHttpClient(defaults)
	rl.Wait()
	return httpclient.Begin().Do("POST", url, headers, body)
}

func (s *Dostow) get(url string, params map[string]string) (*httpclient.Response, error) {
	// Post with files should be sent as multipart.
	rl.Wait()
	return httpclient.Begin().WithHeader("X-DOSTOW-GROUP-ACCESS-KEY", s.key).Get(url, params)
}

func (s *Dostow) put(url string, data string) (*httpclient.Response, error) {
	// Post with files should be sent as multipart.
	headers := make(map[string]string)
	headers["X-DOSTOW-GROUP-ACCESS-KEY"] = s.key
	headers["Content-Type"] = "application/json"
	body := strings.NewReader(data)

	rl.Wait()
	return httpclient.Begin().Do("PUT", url, headers, body)
}

func (s Dostow) Get(id, store string, dst interface{}) (err error) {
	return errors.New("not implemented")
}

func (s Dostow) Save(store string, src interface{}) (key string, err error) {
	var msg map[string]interface{}
	srcjson, _ := json.Marshal(src)
	// logger.Debug("Sending data", "data", string(srcjson))
	resp, err := s.post(s.url+"/store/"+store, string(srcjson))
	// if resp != nil && resp.Body != nil {
	// 	defer resp.Body.Close()
	// }
	// resp, bodyBytes, errs := goreq.New().Post(s.url+"/store/"+store).
	// 	SetHeader("X-DOSTOW-GROUP-ACCESS-KEY", s.key).
	// 	ContentType("json").SendMapString(string(srcjson)).
	// 	SetClient(s.client).
	// 	EndBytes()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case 200:
		if err := json.Unmarshal(bodyBytes, &msg); err == nil {
			return msg["id"].(string), nil
		} else {
			return "", err
		}
	case 404:
		return "", newServerError(resp.StatusCode, string(bodyBytes))
	case 500, 400, 401:
		return "", newServerError(resp.StatusCode, string(bodyBytes))
	default:
		return "", errors.New("Cannot perform action")
	}
}

func (s Dostow) Update(id string, store string, src interface{}) (err error) {
	var msg map[string]interface{}
	srcjson, _ := json.Marshal(src)
	url := s.url + "/store/" + store + "/" + id
	// logger.Debug("Updating data", "url", url, "data", string(srcjson))
	resp, err := s.put(url, string(srcjson))
	// if resp != nil && resp.Body != nil {
	// 	defer resp.Body.Close()
	// }
	// resp, bodyBytes, errs := goreq.New().Put(url).
	// 	SetHeader("X-DOSTOW-GROUP-ACCESS-KEY", s.key).
	// 	ContentType("json").SendMapString(string(srcjson)).
	// 	SetClient(s.client).
	// 	EndBytes()
	if err != nil {
		return handleError(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case 200:
		if err := json.Unmarshal(bodyBytes, &msg); err == nil {
			return nil
		} else {
			return err
		}
	case 404:
		return newServerError(resp.StatusCode, string(bodyBytes))
	case 500, 400, 401:
		return newServerError(resp.StatusCode, string(bodyBytes))
	default:
		return errors.New("Cannot perform action")
	}
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

//FIlterGet gets one item from a store based on some filter
func (s Dostow) FilterGet(filter map[string]interface{}, store string, dst interface{}, opts ObjectStoreOptions) (err error) {
	if filter == nil {
		filter = map[string]interface{}{}
	}
	// request := goreq.New().Get(s.url+"/store/"+store).SetHeader("X-DOSTOW-GROUP-ACCESS-KEY", s.key).SetClient(s.client)
	params := map[string]string{"size": "1"}
	for k, v := range filter {
		// request.QueryData.Add(k, v.(string))
		params[k] = v.(string)
	}
	// request.QueryData.Add("size", "1")
	// resp, bodyBytes, errs := request.EndBytes()

	resp, err := s.get(s.url+"/store/"+store, params)
	// if resp != nil {
	// 	if resp.Body != nil {
	// 		defer resp.Body.Close()
	// 	}
	// }
	if err != nil {
		// logger.Warn("Filter Get Error:")
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	var dl DataList
	switch resp.StatusCode {
	case 200:
		if err := json.Unmarshal(bodyBytes, &dl); err != nil {
			return err
		}
	case 500, 400, 401:
		return newServerError(resp.StatusCode, "Unable to perform action due to server error")
	default:
		return errors.New("Cannot perfrom action")
	}
	if dl.Count >= 1 {

		//		logger.Debug("FilterGet from "+store, "data", dl)
		//		d,_ := json.Marshal(dl.Data[0])
		//		err = json.Unmarshal(d, dst)
		err = mapstructure.Decode(dl.Data[0].(map[string]interface{}), dst)
	} else {
		err = ErrNotFound
	}
	return
}

func (s Dostow) FilterGetAll(filter map[string]interface{}, count int, skip int, store string, opts ObjectStoreOptions) (rrows ObjectRows, err error) {
	params := map[string]string{}
	for k, v := range filter {
		params[k] = to.String(v)
	}
	if count > -1 {
		filter["size"] = string(count)
	}
	// resp, bodyBytes, errs := goreq.New().Get(s.url+"/store/"+store).
	// 	SetHeader("X-DOSTOW-GROUP-ACCESS-KEY", s.key).
	// 	Query(filter).
	// 	SetClient(s.client).
	// 	EndBytes()

	resp, err := s.get(s.url+"/store/"+store, params)
	// if resp != nil && resp.Body != nil {
	// 	defer resp.Body.Close()
	// }

	if err != nil {
		logger.Error("Filter GetAll Error:", "err", err)
		return nil, handleError(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	var dst []interface{}
	switch resp.StatusCode {
	case 200:
		if err := json.Unmarshal(bodyBytes, &dst); err != nil {
			return nil, err
		}
		return DostowRows{dst}, err
	case 500, 400, 401:
		return nil, newServerError(resp.StatusCode, "Unable to perform action due to server error")
	default:
		return nil, errors.New("Cannot perfrom action")
	}
	return nil, errors.New("Not found")
}

func (s Dostow) FilterDelete(filter map[string]interface{}, store string, opts ObjectStoreOptions) (err error) {
	return errors.New("not implemented")
}

func (s Dostow) FilterUpdate(filter map[string]interface{}, src interface{}, store string) (err error) {
	return errors.New("Not Implemented")
}

func (s Dostow) FilterCount(filter map[string]interface{}, store string, opts ObjectStoreOptions) (int64, error) {
	return 0, errors.New("not implemented")
}

func (s Dostow) GetByFieldsByField(name, val, store string, fields []string, dst interface{}) (err error) {
	return errors.New("not implemented")
}

// Streams upload directly from file -> mime/multipart -> pipe -> http-request
func streamingUploadFile(id, field, path, store string, w *io.PipeWriter, file io.Reader) {
	// defer file.Close()
	defer w.Close()
	writer := multipart.NewWriter(w)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		log.Fatal("err", "err", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal("err", "err", err)
		return
	}
	err = writer.Close()
	if err != nil {
		log.Fatal("err", "err", err)
		return
	}
}

// func (s Dostow) UploadFile(id, field, path, store string, file io.Reader) (interface{}, error) {
// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	part, err := writer.CreateFormFile("file", filepath.Base(path))
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, err = io.Copy(part, file)
// 	err = writer.Close()
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, bodyBytes, errs := goreq.New().Post(s.url+"/upload/"+store+"/"+id+"/"+field).
// 		SetHeader("X-DOSTOW-GROUP-ACCESS-KEY", s.key).
// 		SetClient(s.client).
// 		EndBytes()
// 	if errs != nil {
// 		logger.Error("Upload Error:", "err", errs)
// 		return nil, handleError(err)
// 	}
// 	var dst interface{}
// 	switch resp.StatusCode {
// 	case 200:
// 		if err := json.Unmarshal(bodyBytes, &dst); err != nil {
// 			return nil, err
// 		}
// 		return dst, err
// 	case 500, 400, 401:
// 		return nil, newServerError(resp.StatusCode, "Unable to perform action due to server error")
// 	default:
// 		return nil, errors.New("Cannot perfrom action")
// 	}
// 	// return
// }

func (s Dostow) Close() {

}
func handleError(err error) error {
	if err != nil {
		if e, ok := err.(*url.Error); ok {
			//e.Err == *net.OpError
			if e.Err.Error() == "no such host" {
				return newServerError(500, "Server Unreachable")
			}
			return newServerError(500, "Server Unreachable")
		}
		logger.Error("Error", "err", err)
	}
	return newServerError(500, "Unable to complete action")
}

func NewStore(apiUrl, accessKey string) Dostow {
	httpclient.Defaults(httpclient.Map{
		"opt_useragent":   USERAGENT,
		"opt_timeout":     TIMEOUT,
		"Accept-Encoding": "gzip,deflate,sdch",
	})

	s := Dostow{url: apiUrl, key: accessKey, debug: true, client: nil}
	return s
}
