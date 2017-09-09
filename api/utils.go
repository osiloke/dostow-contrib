package api

import (
	"encoding/json"
	"errors"
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
	raw, _, err := source.Store.List(name,
		Query(PaginationParams{Size: 100}),
	)
	if err == nil {
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
	}
	return
}

//CloneStore clones an existing store from one group or server to another
func CloneStore(name, dest string, extraFields map[string]interface{}, scon, dcon *ConnectionParams) {

	source := NewAdminClient(scon.Url, scon.GroupId, scon.AuthorizedToken, http.DefaultClient)
	destination := NewAdminClient(dcon.Url, dcon.GroupId, dcon.AuthorizedToken, http.DefaultClient)

	//create schema
	// log.Info("retrieving schema from source")
	raw, rsp, err := source.Schema.Filter(filter{Name: name})
	// defer log.Trace("copying store").Stop(&err)
	if err == nil {
		if rsp.StatusCode != 200 {
			err = errors.New("could not retrieve schema " + rsp.Status)
			return
		}
		// log.Info("retrieved schema")
		result := new(Result)
		if err = json.Unmarshal(*raw, result); err == nil {
			result.Data[0]["name"] = dest
			schema := &result.Data[0]
			log.Info("creating schema in destination")
			raw, rsp, err = destination.Schema.Create(dest, schema)

			if err != nil {
				if err.Error() == "already exists" {
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
									return
								}
								if raw, rsp, err = destination.Schema.Create(dest, schema); err == nil {
									if rsp.StatusCode != 200 {
										err = errors.New("could not create store in destination " + rsp.Status)
										return
									}
								} else {
									return
								}
							} else {
								return
							}
						}
					}
				}
				// if rsp.StatusCode != 200 {
				// 	err = errors.New("could not create store in destination " + rsp.Status)
				// 	return
				// }
			}
			if err = deleteRows(dest, destination); err == nil {
				copyRows(name, dest, extraFields, source, destination, 0)
			}
		}
	}
	return
}
