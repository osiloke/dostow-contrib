package api

import (
	"encoding/json"

	"github.com/dghubble/sling"
)

type searchQuery struct {
	Q    string `url:"q,omitempty"`
	Size int    `url:"size,omitempty"`
	Skip int64  `url:"skip,omitempty"`
	Sort string `url:"sort,omitempty"`
}

// Opt defines custom options for requests
type Opt func(s *sling.Sling) *sling.Sling

// Authorize authorizes a request with a token
func Authorize(token string) func(s *sling.Sling) *sling.Sling {
	return func(s *sling.Sling) *sling.Sling {
		return s.Set("Authorization", "bearer "+token)
	}
}

// GenericQuery adds a q (query) url param
func GenericQuery(q interface{}) func(s *sling.Sling) *sling.Sling {
	return func(s *sling.Sling) *sling.Sling {
		return s.QueryStruct(q)
	}
}

// Query adds a q (query) url param
func Query(q interface{}) func(s *sling.Sling) *sling.Sling {
	return func(s *sling.Sling) *sling.Sling {
		b, err := json.Marshal(q)
		if err != nil {
			return s.QueryStruct(nil)
		}
		return s.QueryStruct(searchQuery{Q: string(b)})
	}
}

// QueryParams adds a q (query) url param
func QueryParams(q interface{}, size int, skip int64, sort ...interface{}) func(s *sling.Sling) *sling.Sling {
	return func(s *sling.Sling) *sling.Sling {
		query := ""
		if q != nil {
			b, err := json.Marshal(q)
			if err == nil {
				query = string(b)
			}
		}
		if len(sort) > 0 {
			b, err := json.Marshal(sort[0])
			if err == nil {
				return s.QueryStruct(searchQuery{query, size, skip, string(b)})
			}
		}

		if len(query) > 0 {
			return s.QueryStruct(searchQuery{Q: query, Size: size, Skip: skip})
		}
		return s.QueryStruct(searchQuery{Size: size, Skip: skip})
	}
}
