package cos

import (
	"errors"

	"github.com/imroc/req"
)

var (
	ErrForbidden = errors.New("bad athuntication")
	ErrNotFound  = errors.New("not found")
)

const DefaultSignExpireTime = 86400

type Client struct {
	*req.Req
	SecretId  string
	SecretKey string
}

func NewClient(secretId, secretKey string) *Client {
	return &Client{
		Req:       req.New(),
		SecretId:  secretId,
		SecretKey: secretKey,
	}
}

func (client *Client) signRequest(req *request) {
	req.BuildAuth(client.SecretId, client.SecretKey)
	req.Headers["Authorization"] = req.Authorization
}

func (client *Client) buildError(resp *req.Resp) error {
	err := Error{}
	if e := resp.ToXML(&err); e != nil {
		return e
	}
	err.StatusCode = resp.Response().StatusCode
	if err.Message == "" {
		err.Message = resp.Response().Status
	}
	return &err
}

func (client *Client) Do(req *request, options ...Option) (*req.Resp, error) {
	// setup option
	for _, option := range options {
		err := option(req)
		if err != nil {
			return nil, err
		}
	}
	client.signRequest(req) // setup Authorization
	resp, err := client.Req.Do(req.Method, req.Url, req.Headers, req.Params, req.Payload)
	if err != nil {
		return resp, err
	}

	code := resp.Response().StatusCode
	if code != 200 && code != 204 && code != 206 {
		// return nil, fmt.Errorf("bad status: %s", resp.Response().Status)
		return resp, client.buildError(resp)
	}

	return resp, err
}
