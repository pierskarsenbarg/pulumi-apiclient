package pulumiapi

import (
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	c       *http.Client
	token   string
	baseurl url.URL
}

func NewClient(args ...string) *Client {
	c := &http.Client{Timeout: time.Minute}

	var token string

	var baseURL = url.URL{
		Scheme: "https",
		Host:   "api.pulumi.com",
		Path:   "/api/",
	}

	if len(args) == 0 {
		token = ""
	} else {
		token = args[0]
	}

	if len(args) == 2 {
		token = args[0]
		baseURL.Host = args[1]
	}

	return &Client{
		c:       c,
		token:   token,
		baseurl: baseURL,
	}
}
