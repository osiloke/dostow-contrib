package api

import (
	// "fmt"
	// "net/http"

	"github.com/dghubble/sling"
)

type GroupService struct {
	sling *sling.Sling
}

func newGroupService(sling *sling.Sling) *GroupService {
	return &GroupService{
		sling: sling.Path("group/"),
	}
}
