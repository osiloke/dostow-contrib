package api

import (
	"net/http"

	"github.com/dghubble/sling"
	"github.com/ernesto-jimenez/httplogger"
	"log"
	"os"
	"time"
)

type httpLogger struct {
	log *log.Logger
}

func newLogger() *httpLogger {
	return &httpLogger{
		log: log.New(os.Stderr, "dostow - ", log.LstdFlags),
	}
}

func (l *httpLogger) LogRequest(req *http.Request) {
	l.log.Printf(
		"Request %s %s",
		req.Method,
		req.URL.String(),
	)
	l.log.Printf(
		"Headers %v",
		req.Header,
	)
}

func (l *httpLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	duration /= time.Millisecond
	if err != nil {
		l.log.Println(err)
	} else {
		l.log.Printf(
			"Response method=%s status=%d durationMs=%d %s",
			req.Method,
			res.StatusCode,
			duration,
			req.URL.String(),
		)
	}
}

// Client is a Dostow client for making Dostow API requests.
type Client struct {
	sling  *sling.Sling
	Group  *GroupService
	Schema *SchemaService
	Auth   *AuthService
	Store  *StoreService
	File   *FileService
}

// NewClient return a new Client
func NewClient(apiUrl, apiKey string, httpClients ...*http.Client) *Client {
	var httpClient *http.Client

	if len(httpClients) > 0 {
		httpClient = httpClients[0]
	} else {
		httpClient = &http.Client{
			Transport: httplogger.NewLoggedTransport(http.DefaultTransport, newLogger()),
		}
	}

	base := sling.New().Client(httpClient).Base(apiUrl).Set("X-Dostow-Group-Access-Key", apiKey)
	return &Client{
		sling:  base,
		Schema: newSchemaService(base.New()),
		Auth:   newAuthService(base.New()),
		Store:  newStoreService(base.New()),
	}
}

func NewAdminClient(apiUrl, groupId, token string, httpClients ...*http.Client) *Client {
	var httpClient *http.Client

	if len(httpClients) > 0 {
		httpClient = httpClients[0]
	} else {
		httpClient = http.DefaultClient
	}
	base := sling.New().Client(httpClient).Base(apiUrl).Set("X-Dostow-Group", groupId).Set("Authorization", token)
	return &Client{
		sling:  base,
		Group:  newGroupService(base.New()),
		Schema: newSchemaService(base.New()),
		Auth:   newAuthService(base.New()),
		Store:  newStoreService(base.New()),
		File:   newFileService(base.New()),
	}
}
