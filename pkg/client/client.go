package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"callback-server/pkg/cerr"
	"callback-server/pkg/header"
)

type Client struct {
	addr       string
	httpClient *http.Client
}

type NewOpt func(client *Client)

func NewClient(addr string, opts ...NewOpt) *Client {
	c := &Client{
		addr:       addr + "/cb/",
		httpClient: http.DefaultClient,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func WithHttpClient(hc *http.Client) NewOpt {
	return func(c *Client) {
		c.httpClient = hc
	}
}

type Data struct {
	ContentType string
	Body        []byte
}

func (c *Client) Call(ctx context.Context, id string, data *Data) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.addr+id,
		bytes.NewReader(data.Body))
	if err != nil {
		return err
	}
	req.Header.Set(header.ContentType, data.ContentType)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return cerr.HttpStatusError{
			Code: resp.StatusCode,
			Body: b,
		}
	}
	return nil
}

func (c *Client) Wait(ctx context.Context, id string) (*Data, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.addr+id, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, cerr.HttpStatusError{
			Code: resp.StatusCode,
			Body: b,
		}
	}
	return &Data{
		ContentType: resp.Header.Get(header.ContentType),
		Body:        b,
	}, nil
}
