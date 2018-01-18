package api

import (
	"encoding/json"
	"github.com/dghubble/sling"
)

type searchQuery struct {
	Q string `url:"q"`
}

// Opt defines custom options for requests
type Opt func(s *sling.Sling) *sling.Sling

// Authorize authorizes a request with a token
func Authorize(token string) func(s *sling.Sling) *sling.Sling {
	return func(s *sling.Sling) *sling.Sling {
		return s.Set("Authorization", "bearer "+token)
	}
}

// Query adds a q (query) url param
func Query(q interface{}) func(s *sling.Sling) *sling.Sling {
	return func(s *sling.Sling) *sling.Sling {
		b, err := json.Marshal(q)
		if err != nil {
			return s.QueryStruct(nil)
		}
		return s.QueryStruct(searchQuery{string(b)})
	}
}
