package cos

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/imroc/req"
)

type Bucket struct {
	*Client
	url       string
	secretId  string
	secretKey string
}

func NewBucketFromURL(url string, client *Client) *Bucket {
	return &Bucket{
		Client: client,
		url:    strings.TrimRight(url, "/"),
	}
}

func (b *Bucket) do(method, path string, options ...Option) (*req.Resp, error) {
	req, err := b.newRequest(method, path)
	if err != nil {
		return nil, err
	}
	return b.Do(req, options...)
}

func (b *Bucket) newRequest(method, path string) (*request, error) {
	url := strings.Join([]string{b.url, path}, "")
	return NewRequest(method, url)
}

func (b *Bucket) Exists() (bool, error) {
	req, err := b.newRequest("HEAD", "/")
	if err != nil {
		return false, err
	}
	_, err = b.Do(req)
	if err != nil {
		return false, err
	}
	return true, nil
}

// The ListObjectsResult type holds the results of a List bucket operation.
type ListObjectsResult struct {
	XMLName xml.Name `xml:"ListBucketResult"`
	Name    string   `xml:"Name"`
	Prefix  string   `xml:"Prefix"`
	Marker  string   `xml:"Marker"`
	MaxKeys int      `xml:"MaxKeys"`
	// IsTruncated is true if the results have been truncated because
	// there are more keys and prefixes than can fit in MaxKeys.
	// N.B. this is the opposite sense to that documented (incorrectly) in
	// http://goo.gl/YjQTc
	IsTruncated    bool               `xml:"IsTruncated"`
	Objects        []ObjectProperties `xml:"Contents"`
	CommonPrefixes []string           `xml:"CommonPrefixes>Prefix"`
	// if IsTruncated is true, pass NextMarker as marker argument to List()
	// to get the next set of keys
	NextMarker string `xml:"NextMarker"`
}

// The ObjectProperties type represents an item stored in an bucket.
type ObjectProperties struct {
	XMLName      xml.Name  `xml:"Contents"`
	Key          string    `xml:"Key"`          // Object key
	LastModified time.Time `xml:"LastModified"` // Object last modified time
	Size         int64     `xml:"Size"`         // Object size
	// ETag gives the hex-encoded MD5 sum of the contents,
	// surrounded with double-quotes.
	ETag         string `xml:"ETag"`         // Object ETag
	Owner        Owner  `xml:"Owner"`        // Object owner information
	StorageClass string `xml:"StorageClass"` // Object storage class (Standard, IA, Archive)
}

// The Owner type represents the owner of the object in an bucket.
type Owner struct {
	XMLName     xml.Name `xml:"Owner"`
	ID          string   `xml:"ID"`
	DisplayName string   `xml:"DisplayName"`
}

func (b *Bucket) ListObjects(options ...Option) (*ListObjectsResult, error) {
	resp, err := b.do("GET", "/", options...)
	if err != nil {
		return nil, err
	}
	var result ListObjectsResult
	err = resp.ToXML(&result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

// GetObjectResult is the result of DoGetObject
type GetObjectResult struct {
	Response *http.Response
}

func (b *Bucket) DoGetObject(key string, options ...Option) (*GetObjectResult, error) {
	resp, err := b.do("GET", key, options...)
	if err != nil {
		return nil, err
	}
	return &GetObjectResult{resp.Response()}, nil
}

func (b *Bucket) GetObject(key string, options ...Option) (io.ReadCloser, error) {
	result, err := b.DoGetObject(key, options...)
	if err != nil {
		return nil, err
	}
	return result.Response.Body, nil
}

func (b *Bucket) GetObjectMeta(key string) (http.Header, error) {
	resp, err := b.do("HEAD", key)
	if err != nil {
		return nil, err
	}
	return resp.Response().Header, nil
}

func (b *Bucket) DeleteObject(key string) error {
	_, err := b.do("DELETE", key)
	if err != nil {
		return err
	}
	return nil
}

// PutObject upload a file, body could be io.Reader, []byte, string
func (b *Bucket) PutObject(key string, body interface{}, options ...Option) error {
	_, err := b.do("PUT", key, Body(body), ContentType(getContentType(key)))
	if err != nil {
		return err
	}
	return nil
}
