package api

// PaginationParams ...
type PaginationParams struct {
	Before int `url:"before,omitempty"`
	After  int `url:"after,omitempty"`
	Size   int `url:"size,omitempty"`
}

// StatusDeletion ...
type StatusDeletion struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}
