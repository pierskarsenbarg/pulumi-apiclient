package pulumiapi

import (
	"net/http"
	"net/url"
	"time"
)

var baseURL = url.URL{
	Scheme: "https",
	Host:   "api.pulumi.com",
	Path:   "/api/",
}

type Client struct {
	c     *http.Client
	token string
}

func NewClient(token string) *Client {
	c := &http.Client{Timeout: time.Minute}

	return &Client{
		c:     c,
		token: token,
	}
}
