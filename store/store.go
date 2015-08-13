package store

import (
	"github.com/mgutz/logxi/v1"
	. "github.com/osiloke/gostore"
	"github.com/bndr/gopencils"
	"github.com/dustin/gojson"
	"errors"
	"net/url"
)

var logger = log.New("gostore.dostow")

type ServerError struct{
	code int
	msg string
}

func (e ServerError) Error() string{
	return e.msg
}

func newServerError(code int, msg string) ServerError{
	return ServerError{code, msg}
}

func NewStore(apiUrl, accessKey string) Dostow {
	s := Dostow{gopencils.Api(apiUrl+"/store"), accessKey}
	return s
}

type Dostow struct {
	api *gopencils.Resource
	key string
}

type DostowRows struct {
}

func (s Dostow) req(store string) *gopencils.Resource{
	resource := s.api.Res(store)
	resource.SetHeader("X-DOSTOW-GROUP-ACCESS-KEY", s.key)
	return resource
}

func (s DostowRows) Next(dst interface{}) (bool, error) {

	return false, nil
}

func (s DostowRows) Close(){
}

func (s Dostow) CreateDatabase() (err error) {
	return errors.New("not implemented")
}

func (s Dostow) GetStore() interface{} {
	return s.api
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

func (s Dostow) FilterBefore(id string, filter map[string]interface{}, count int, skip int, store string) (rows ObjectRows, err error){
	return nil, errors.New("not implemented")
}

func (s Dostow) FilterBeforeCount(id string, filter map[string]interface{}, count int, skip int, store string) (int64, error){
	return 0, errors.New("not implemented")
}

func (s Dostow) FilterSince(id string, filter map[string]interface{}, count int, skip int, store string) (rows ObjectRows, err error){
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
	r, err := s.req(store).Post(src)
	if err != nil {
//		log.Info("err", "err", err)
		if e, ok := err.(*url.Error); ok{
			//e.Err == *net.OpError
			if e.Err.Error() == "no such host"{
				return "", newServerError(500, "Server Unreachable")
			}
			return "", newServerError(500, "Server Unreachable")
		}
		return "", newServerError(500, "Unable to complete taction")

	}

	decoder := json.NewDecoder(r.Raw.Body)
	decoder.Decode(&msg)
	switch(r.Raw.StatusCode){
	case 200:
		return msg["id"].(string), nil
	case 500,400,401:
		return "", newServerError(r.Raw.StatusCode, "Unable to perform action due to server error")
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
	return errors.New("not implemented")
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
