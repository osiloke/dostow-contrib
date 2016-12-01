package api

import (
	"github.com/dghubble/sling"
)

type Opt func(sling *sling.Sling)

func Authorize(token string) func(sling *sling.Sling) {
	return func(sling *sling.Sling) {
		sling.Set("Authorization", "bearer "+token)
	}
}
