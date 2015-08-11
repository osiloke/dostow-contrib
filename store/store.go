package store

import (
	"github.com/mgutz/logxi/v1"
	"strings"
	. "github.com/osiloke/gostore"
	"fmt"
	"github.com/dustin/gojson"
)

var logger = log.New("gostore.rethink")

func NewStore(session *r.Session, database string) Dostow {
	s := Dostow{session, database}
	s.CreateDatabase()
	return s
}

type Dostow struct {
	Session  *r.Session
	Database string
}

type DostowRows struct {
	cursor *r.Cursor
}

func (s DostowRows) Next(dst interface{}) (bool, error) {
	if !s.cursor.Next(dst) {
		//		logger.Debug("Error getting next", "err", s.cursor.Err(), "isNil", s.cursor.IsNil())
		return false, s.cursor.Err()
	}
	return true, nil
}

func (s DostowRows) Close(){
	s.cursor.Close()
}

func (s Dostow) CreateDatabase() (err error) {
	return r.DBCreate(s.Database).Exec(s.Session)
}

func (s Dostow) GetStore() interface{} {
	return s.Session
}

func (s Dostow) CreateTable(store string, sample interface{}) (err error) {
	_, err = r.DB(s.Database).TableCreate(store).RunWrite(s.Session)

	return
}

func (s Dostow) All(count int, skip int, store string) (rrows ObjectRows, err error) {
	result, err := r.DB(s.Database).Table(store).OrderBy(r.OrderByOpts{Index:r.Desc("id")}).Run(s.Session)
	if err != nil {
		return
	}
	rrows = DostowRows{result}
	return
}

func (s Dostow) AllCursor(store string) (ObjectRows, error) {
	result, err := r.DB(s.Database).Table(store).Run(s.Session)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	return DostowRows{result}, nil
}

//Before will retrieve all old rows that were created before the row with id was created
// [1, 2, 3, 4], before 2 will return [3, 4]
//r.db('worksmart').table('store').orderBy({'index': r.desc('id')}).filter(r.row('schemas')
// .eq('osiloke_tsaboin_silverbird').and(r.row('id').lt('55b54e93f112a16514000057')))
// .pluck('schemas', 'id','tid', 'timestamp', 'created_at').limit(100)
func (s Dostow) Before(id string, count int, skip int, store string) (rows ObjectRows, err error) {
	result, err := r.DB(s.Database).Table(store).Filter(r.Row.Field("id").Lt(id)).Limit(count).Skip(skip).Run(s.Session)
	if err != nil {
		return
	}
	defer result.Close()
	//	result.All(dst)
	rows = DostowRows{result}
	return
}

func (s Dostow) FilterBefore(id string, filter map[string]interface{}, count int, skip int, store string) (rows ObjectRows, err error){
	result, err := r.DB(s.Database).Table(store).Between(
		r.MinVal, id, r.BetweenOpts{RightBound: "closed"}).OrderBy(
		r.OrderByOpts{Index:r.Desc("id")}).Filter(
		filter).Limit(count).Run(s.Session)
	if err != nil {
		return
	}
	//	var dst interface{}
	f, _ := json.Marshal(filter)
	logger.Debug("FilterBefore", "query",
		fmt.Sprintf("r.db('%s').table('%s').between(r.minval, '%s').orderBy({index:r.desc('id')}).filter(%s).limit(%d)",
			s.Database, store, id, string(f), count))
	rows = DostowRows{result}
	return
}

func (s Dostow) FilterBeforeCount(id string, filter map[string]interface{}, count int, skip int, store string) (int64, error){
	result, err := r.DB(s.Database).Table(store).Between(
		r.MinVal, id).OrderBy(
		r.OrderByOpts{Index:r.Desc("id")}).Filter(
		filter).Count().Run(s.Session)
	defer result.Close()

	var cnt int64
	if err = result.One(&cnt); err != nil {
		return 0, ErrNotFound
	}
	return cnt, nil
}

func (s Dostow) FilterSince(id string, filter map[string]interface{}, count int, skip int, store string) (rows ObjectRows, err error){
	result, err := r.DB(s.Database).Table(store).Between(
		id, r.MaxVal, r.BetweenOpts{LeftBound: "open", Index:"id"}).OrderBy(
		r.OrderByOpts{Index:r.Desc("id")}).Filter(
		filter).Limit(count).Run(s.Session)
	if err != nil {
		return
	}
	//	result.All(dst)
	rows = DostowRows{result}
	return
}

//This will retrieve all new rows that were created since the row with id was created
// [1, 2, 3, 4], since 2 will return [1]
func (s Dostow) Since(id string, count, skip int, store string) (rrows ObjectRows, err error) {
	result, err := r.DB(s.Database).Table(store).Filter(r.Row.Field("id").Gt(id)).Limit(count).Skip(skip).Run(s.Session)
	if err != nil {
		return
	}
	//	result.All(dst)
	rrows = DostowRows{result}
	return
}

func (s Dostow) Get(id, store string, dst interface{}) (err error) {
	result, err := r.DB(s.Database).Table(store).Get(id).Run(s.Session)
	if err != nil {
		//		logger.Error("Get", "err", err)
		return
	}
	defer result.Close()
	if result.Err() != nil {
		return result.Err()
	}
	if err = result.One(dst); err == r.ErrEmptyResult {
		//		logger.Error("Get", "err", err)
		return ErrNotFound
	}
	return nil
}

func (s Dostow) Save(store string, src interface{}) (key string, err error) {
	result, err := r.DB(s.Database).Table(store).Insert(src, r.InsertOpts{Durability: "soft"}).RunWrite(s.Session)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate primary key") {
			err = ErrDuplicatePk
		}
		return
	}
	if len(result.GeneratedKeys) > 0 {
		key = result.GeneratedKeys[0]
	}
	return

}

func (s Dostow) Update(id string, store string, src interface{}) (err error) {
	_, err = r.DB(s.Database).Table(store).Get(id).Update(src, r.UpdateOpts{Durability: "soft"}).RunWrite(s.Session)
	return

}

func (s Dostow) Delete(id string, store string) (err error) {
	_, err = r.DB(s.Database).Table(store).Get(id).Delete(r.DeleteOpts{Durability: "hard"}).RunWrite(s.Session)
	return
}

func (s Dostow) Stats(store string) (map[string]interface{}, error) {
	result, err := r.DB(s.Database).Table(store).Count().Run(s.Session)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	var cnt int64
	if err = result.One(&cnt); err != nil {
		return nil, ErrNotFound
	}
	return map[string]interface{}{"count": cnt}, nil
}

func (s Dostow) GetByField(name, val, store string, dst interface{}) (err error) {
	result, err := r.DB(s.Database).Table(store).Filter(r.Row.Field(name).Eq(val)).Run(s.Session)
	if err != nil {
		return
	}
	defer result.Close()
	if err = result.One(dst); err == r.ErrEmptyResult {
		return ErrNotFound
	}
	return
}

func (s Dostow) FilterGet(filter map[string]interface{}, store string, dst interface{}) (err error) {
	result, err := r.DB(s.Database).Table(store).Filter(filter).Run(s.Session)
	if err != nil {
		return
	}
	defer result.Close()
	if err = result.One(dst); err == r.ErrEmptyResult {
		return ErrNotFound
	}
	return
}

func (s Dostow) FilterGetAll(filter map[string]interface{}, count int, skip int, store string) (rrows ObjectRows, err error) {
	//	logger.Debug("Filter get all", "store", store, "filter", filter)
	result, err := r.DB(s.Database).Table(store).OrderBy(
		r.OrderByOpts{Index:r.Desc("id")}).Filter(filter).Limit(count).Skip(skip).Run(s.Session)
	if err != nil {
		return
	}
	rrows = DostowRows{result}
	return
}

func (s Dostow) FilterDelete(filter map[string]interface{}, store string) (err error) {
	_, err = r.DB(s.Database).Table(store).Filter(filter).Delete(r.DeleteOpts{Durability: "soft"}).RunWrite(s.Session)
	if err == r.ErrEmptyResult {
		return ErrNotFound
	}
	return
}

func (s Dostow) FilterCount(filter map[string]interface{}, store string) (int64, error) {
	result, err := r.DB(s.Database).Table(store).Filter(filter).Count().Run(s.Session)
	if err != nil {
		return 0, err
	}
	defer result.Close()
	var cnt int64
	if err = result.One(&cnt); err != nil {
		return 0, ErrNotFound
	}
	return cnt, nil
}

func (s Dostow) GetByFieldsByField(name, val, store string, fields []string, dst interface{}) (err error) {
	result, err := r.DB(s.Database).Table(store).Filter(r.Row.Field(name).Eq(val)).Pluck(fields).Run(s.Session)
	if err != nil {
		return
	}
	defer result.Close()
	if err = result.One(dst); err == r.ErrEmptyResult {
		return ErrNotFound
	}
	return
}

func (s Dostow) Close() {
	s.Session.Close()
}
