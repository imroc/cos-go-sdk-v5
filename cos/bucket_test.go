package cos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/k0kubun/pp"

	"github.com/imroc/req"
)

var bucket *Bucket

func init() {
	req.Debug = true
	bs, err := ioutil.ReadFile("../testdata/bucket.json")
	if err != nil {
		panic(err)
	}
	var Config struct {
		Url       string `json:"url"`
		SecretId  string `json:"secretId"`
		SecretKey string `json:"secretKey"`
	}
	err = json.Unmarshal(bs, &Config)
	if err != nil {
		panic(err)
	}
	client := NewClient(Config.SecretId, Config.SecretKey)
	bucket = NewBucketFromURL(Config.Url, client)
}

func TestBucketExists(t *testing.T) {
	exists, err := bucket.Exists()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("exists: %v", exists)
}

func TestBucketListObjects(t *testing.T) {
	result, err := bucket.ListObjects(Prefix("/test/"), Delimiter(""), MaxKeys(30))
	// listResp, err := bucket.ListObjects(Delimiter("/"), MaxKeys(30))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("listResp: %+v", *result)
}

func TestBucketGetObjects(t *testing.T) {
	rc, err := bucket.GetObject("/test/test.png")
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("test.png", bs, 0666)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBucketPutObjects(t *testing.T) {
	file, err := os.Open("util.go")
	if err != nil {
		t.Fatal(err)
	}
	err = bucket.PutObject("/test/util.go", file)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBucketDeleteObjects(t *testing.T) {
	err := bucket.DeleteObject("/test/LICENSE")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBucketGetObjectsMeta(t *testing.T) {
	header, err := bucket.GetObjectMeta("/test/test.png")
	if err != nil {
		t.Fatal(err)
	}
	pp.Println(header)
}
func TestSomething(t *testing.T) {
	typ := getContentType("/test/dsklfjs哈哈adffukc嘻嘻.xml")
	fmt.Println("type:", typ)
}
