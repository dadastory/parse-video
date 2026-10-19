package parser

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/wujunwei928/parse-video/utils"
)

// ErrNoMedia means a page did not supply a usable note or downloadable media.
var ErrNoMedia = errors.New("page contains no downloadable media")

// Client is an isolated library parser. It does not consult global proxy settings.
// Configure proxy, timeout, and transport on the supplied standard HTTP client.
// The client must not be modified after construction and is safe for concurrent use.
type Client struct{ httpClient http.Client }

func NewClient(client *http.Client) *Client {
	if client == nil {
		client = &http.Client{}
	}
	return &Client{httpClient: *client}
}

type requestClient struct{ factory func() *resty.Client }

func (c requestClient) newClient() *resty.Client {
	if c.factory != nil {
		return c.factory()
	}
	return newClient()
}

// Parse accepts either a sharing URL or the original surrounding sharing text.
// Platform receivers share this invocation's context, without global mutations.
func (c *Client) Parse(ctx context.Context, text string) (*VideoParseInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	source, err := utils.RegexpMatchUrlFromString(text)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(source)
	if err != nil || u.Hostname() == "" || u.User != nil || u.Port() != "" || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, errors.New("invalid sharing URL")
	}
	client := requestClient{factory: func() *resty.Client {
		native := c.httpClient
		return resty.NewWithClient(&native).OnBeforeRequest(func(_ *resty.Client, req *resty.Request) error { req.SetContext(ctx); return ctx.Err() })
	}}
	for _, info := range sourceInfoMapping(client) {
		for _, host := range info.VideoShareUrlDomain {
			if strings.EqualFold(strings.TrimPrefix(strings.ToLower(u.Hostname()), "www."), strings.TrimPrefix(host, "www.")) {
				result, err := info.VideoShareUrlParser.parseShareUrl(u.String())
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				return result, err
			}
		}
	}
	return nil, errors.New("unsupported sharing source")
}
