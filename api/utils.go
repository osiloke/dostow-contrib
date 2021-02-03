package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/apex/log"
)

func deleteRows(name string, client *Client) (err error) {
	defer log.Trace("copying rows").Stop(&err)
	var size int64
	for {
		_, _, err = client.Store.Clear(name)
		if err != nil {
			if err.Error() == "does not exist or already cleared" {
				err = nil
			}
			break
		}
		size, err = client.Store.Size(name)
		if size == 0 {
			break
		}
	}
	return
}

func copyRows(name, dest string, extraFields map[string]interface{}, source, destination *Client, skip int) (err error) {
	defer log.WithFields(log.Fields{
		"skip": skip,
	}).Trace("copying rows").Stop(&err)
	log.Infof("copying rows from %s to %s", name, dest)
	var raw *json.RawMessage
	raw, _, err = source.Store.List(name,
		Query(PaginationParams{Size: 100}),
	)
	if err != nil {
		if !strings.Contains(err.Error(), "cleared") {
			return err
		}
	}
	rows := new(Result)
	if err = json.Unmarshal(*raw, rows); err == nil {
		cnt := len(rows.Data) + skip
		if cnt <= rows.Total {
			copyRows(name, dest, extraFields, source, destination, cnt)
		} else {
			for _, row := range rows.Data {
				for k, v := range extraFields {
					row[k] = v
				}
				// if raw, err := destination.Store.Create(dest, &row); err == nil {
				// 	if err = json.Unmarshal(*raw, &row); err != nil {
				// 		log.WithFields(log.Fields{"row": row}).Info("created row")
				// 	}
				// }
			}
			if _, err := destination.Store.BulkCreate(dest, rows.Data); err == nil {
				log.Info("created rows")
			}
		}
	}
	return
}

func copySchema(name, dest string, source, destination *Client, copyChain string) error {
	raw, rsp, err := source.Schema.Filter(filter{Q: fmt.Sprintf(`{"name":"%s"}`, name)})
	// defer log.Trace("copy schema").Stop(&err)
	if err != nil {
		return fmt.Errorf("could not retrieve schema %v", err)
	}
	if rsp.StatusCode != 200 {
		return errors.New("could not retrieve schema " + rsp.Status)
	}
	copyChain = copyChain + " => " + name
	// log.Info("retrieved schema")
	result := new(Result)
	if err = json.Unmarshal(*raw, result); err == nil {
		result.Data[0]["name"] = dest
		schema := result.Data[0]
		log.Infof("%s => creating %s at destination", copyChain, dest)
		raw, rsp, err = destination.Schema.Create(dest, &schema)

		if err != nil {
			log.Warnf("received error %v", err)
			if strings.Contains(err.Error(), "conditions not met") {
				log.Warn("backup store already exists")
				err = nil
			} else if e, ok := err.(*APIError); ok {
				if strings.Contains(e.Message, "`") && strings.Contains(e.Message, "cannot be updated") && strings.Contains(e.Message, "does not exist") {
					parts := strings.Split(e.Message, "`")
					if parts[1] == "static" {
						log.Info("creating " + parts[1] + " referenced store")
						if raw, rsp, err = destination.Schema.Create(parts[1], map[string]interface{}{"name": "static"}); err == nil {
							if rsp.StatusCode != 200 {
								err = errors.New("could not create referenced store " + rsp.Status)

							} else {
								if raw, rsp, err = destination.Schema.Create(dest, &schema); err == nil {
									if rsp.StatusCode != 200 {
										err = errors.New("could not create store in destination " + rsp.Status)
									}
								}
							}
						}
					}
				} else if strings.Contains(e.Message, "`") && strings.Contains(e.Message, "does not exist") {
					// try to create store
					// var schemaMap map[string]interface{}
					// json.Unmarshal([]byte(schema), &schemaMap)
					if strings.HasPrefix(e.Message, "ref") {
						vals := strings.Split(e.Message, "`")
						fieldName := vals[1]
						linkedStore := vals[3]
						log.Infof("%s references %s store", fieldName, linkedStore)
						if err = copySchema(linkedStore, linkedStore, source, destination, copyChain); err == nil {
							return copySchema(name, dest, source, destination, copyChain)
						}
					} else if strings.HasPrefix(e.Message, "cannot include") {
						vals := strings.Split(e.Message, "`")
						includedField := vals[1]
						includedStore := vals[3]
						log.Warnf("%s includes link `%s` => %s store", name, includedField, includedStore)
						// is the included item a store we are already working on
						if strings.Contains(copyChain, includedStore) {
							// remove include and recreate and then update
							include := schema["include"].(map[string]interface{})
							includeConflict := include[includedField].(map[string]interface{})
							delete(include, includedField)
							schema["include"] = include
							log.Infof("creating %s without %s", name, includedField)
							raw, rsp, err = destination.Schema.Create(dest, &schema)
							if err == nil {
								json.Unmarshal(*raw, &schema)
								schemaID := schema["id"].(string)
								log.Infof("updating %s with include.%s = %v", name, includedField, includeConflict)
								destination.Schema.Update(schemaID, &map[string]interface{}{
									"include": map[string]interface{}{
										includedField: includeConflict,
									},
								})
							}
						}
					}

				}
			}
		}
	}
	return err
}

//CopySchema copy an existing schema from one group or server to another
func CopySchema(name, dest string, scon, dcon *ConnectionParams) {

	source := NewAdminClient(scon.Url, scon.GroupId, scon.AuthorizedToken, http.DefaultClient)
	destination := NewAdminClient(dcon.Url, dcon.GroupId, dcon.AuthorizedToken, http.DefaultClient)

	//create schema
	// log.Info("retrieving schema from source")
	err := copySchema(name, dest, source, destination, "")
	defer log.Trace("copying store").Stop(&err)

	return
}

//CloneStore clones an existing store from one group or server to another
func CloneStore(name, dest string, extraFields map[string]interface{}, scon, dcon *ConnectionParams) {

	source := NewAdminClient(scon.Url, scon.GroupId, scon.AuthorizedToken, http.DefaultClient)
	destination := NewAdminClient(dcon.Url, dcon.GroupId, dcon.AuthorizedToken, http.DefaultClient)

	//create schema
	log.Infof("retrieving store from source %v", scon.Url)
	err := copySchema(name, dest, source, destination, "")
	defer log.Trace("CloneStore").Stop(&err)

	if err == nil {
		err = deleteRows(dest, destination)
		if err != nil {
			if !strings.Contains(err.Error(), "cleared") {
				return
			}
		}
		copyRows(name, dest, extraFields, source, destination, 0)

	}
	return
}
