package api

//ConnectionParams used to define an api connection
type ConnectionParams struct {
	Url             string
	GroupId         string
	AuthorizedToken string
}

//Result defines a list result
type Result struct {
	Data  []map[string]interface{} `json:"data"`
	Total int                      `json:"total_count"`
}

// PaginationParams ...
type PaginationParams struct {
	Before int   `url:"before,omitempty" json:"before,omitempty"`
	After  int   `url:"after,omitempty" json:"after,omitempty"`
	Size   int   `url:"size,omitempty" json:"size,omitempty"`
	Skip   int64 `url:"skip,omitempty" json:"skip,omitempty"`
}

// StatusDeletion ...
type StatusDeletion struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

type filter struct {
	Q string `url:"q"`
}
