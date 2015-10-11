package store

import (
	"errors"
	"net/url"
	"github.com/mgutz/logxi/v1"
	"github.com/mitchellh/mapstructure"
	. "github.com/osiloke/gostore"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
)

var logger = log.New("gostore.dostow")

type ServerError struct {
	code int
	msg  string
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
	url   string
	key   string
	req *gorequest.SuperAgent
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
	srcjson, _ := json.Marshal(src)
	logger.Debug("Sending data", "data", string(srcjson))
	resp, bodyBytes, errs := s.req.Post(s.url + "/store/" + store).
	Set("X-DOSTOW-GROUP-ACCESS-KEY", s.key).
	SendString(string(srcjson)).
	EndBytes()
	if errs != nil {
		return "", handleError(err)
	}
	switch resp.StatusCode {
	case 200:
		if err := json.Unmarshal(bodyBytes, &msg); err == nil {
			return msg["id"].(string), nil
		}else{
			return "", err
		}
	case 404:
		return "", newServerError(resp.StatusCode, string(bodyBytes))
	case 500, 400, 401:
		return "", newServerError(resp.StatusCode, string(bodyBytes))
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

//FIlterGet gets one item from a store based on some filter
func (s Dostow) FilterGet(filter map[string]interface{}, store string, dst interface{}) (err error) {
	if filter == nil {
		filter = map[string]interface{}{}
	}
	request := s.req.Get(s.url + "/store/" + store).Set("X-DOSTOW-GROUP-ACCESS-KEY", s.key)
	for k, v := range filter{
		request.QueryData.Add(k, v.(string))
	}
	request.QueryData.Add("size", "1")
	resp, bodyBytes, errs := request.EndBytes()
	if errs != nil {
		logger.Error("Filter Get Error:", "err", errs)
		return handleError(err)
	}

	var dl DataList
	switch resp.StatusCode {
	case 200:
		if err := json.Unmarshal(bodyBytes, &dl); err != nil{
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
	}else{
		err = ErrNotFound
	}
	return
}

func (s Dostow) FilterGetAll(filter map[string]interface{}, count int, skip int, store string) (rrows ObjectRows, err error) {
	if filter == nil {
		filter = map[string]interface{}{}
	}
	if count > -1 {
		filter["size"] = string(count)
	}
	resp, bodyBytes, errs := s.req.Get(s.url + "/store/" + store).
	Set("X-DOSTOW-GROUP-ACCESS-KEY", s.key).
	Query(filter).
	EndBytes()
	if errs != nil {
		logger.Error("Filter GetAll Error:", "err", errs)
		return nil, handleError(err)
	}
	var dst []interface{}
	switch resp.StatusCode {
	case 200:
		if err := json.Unmarshal(bodyBytes, &dst); err != nil{
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
func handleError(err error) error {
	if e, ok := err.(*url.Error); ok {
		//e.Err == *net.OpError
		if e.Err.Error() == "no such host" {
			return newServerError(500, "Server Unreachable")
		}
		return newServerError(500, "Server Unreachable")
	}
	logger.Error("Error", "err", err)
	return newServerError(500, "Unable to complete action")
}


func NewStore(apiUrl, accessKey string) Dostow {
	s := Dostow{url: apiUrl, key: accessKey, debug: true, req: gorequest.New()}
	return s
}
