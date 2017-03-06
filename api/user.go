package api

import (
	// "fmt"
	// "net/http"
	"encoding/json"
	"github.com/dghubble/sling"
)

type AuthService struct {
	sling *sling.Sling
}

type SignInRequest struct {
	Email    string `json:"email"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginSuccess struct {
	Data     map[string]interface{} `json:"data"`
	Token    string                 `json:"token"`
	Username string                 `json:"username"`
}

func newAuthService(sling *sling.Sling) *AuthService {
	return &AuthService{
		sling: sling.Path("auth/"),
	}
}

func (s *AuthService) Me(token string, user interface{}) error {
	apiError := new(APIError)
	_s := s.sling.New().Get("me").Set("Authorization", "bearer "+token)
	_, err := _s.Receive(user, apiError)
	return relevantError(err, apiError)
}

func (s *AuthService) SignIn(req *SignInRequest) (*LoginSuccess, error) {
	var result *LoginSuccess = &LoginSuccess{}
	apiError := new(APIError)
	_s := s.sling.New().Post("sign_in").BodyJSON(req)
	_, err := _s.Receive(result, apiError)
	return result, relevantError(err, apiError)
}

func (s *AuthService) Register(req *RegisterRequest) (*json.RawMessage, error) {
	var result *json.RawMessage = &json.RawMessage{}
	apiError := new(APIError)
	_s := s.sling.New().Post("register").BodyJSON(req)
	_, err := _s.Receive(result, apiError)
	return result, relevantError(err, apiError)
}

func (s *AuthService) ResetPassword(req *ResetPasswordRequest) (*json.RawMessage, error) {
	var result *json.RawMessage = &json.RawMessage{}
	apiError := new(APIError)
	_s := s.sling.New().Post("reset_password").BodyJSON(req)
	_, err := _s.Receive(result, apiError)
	return result, relevantError(err, apiError)
}
