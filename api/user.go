package api

import (
	// "fmt"
	// "net/http"

	"github.com/dghubble/sling"
)

type AuthService struct {
	sling *sling.Sling
}

func newAuthService(sling *sling.Sling) *AuthService {
	return &AuthService{
		sling: sling.Path("auth/"),
	}
}
