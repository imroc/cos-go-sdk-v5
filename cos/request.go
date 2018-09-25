package cos

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/imroc/req"
)

type request struct {
	Method   string
	Url      string
	Params   req.QueryParam
	Headers  req.Header
	prepared bool
	Payload  interface{}
	timeout  time.Duration
	expire   time.Time //time point

	SignKey       string
	FormatString  string
	StringToSign  string
	Authorization string

	SignTime     string
	KeyTime      string
	HeaderList   string
	UrlParamList string
	Signature    string
}

func NewRequest(method, urlStr string) (*request, error) {
	req := &request{
		Method:  method,
		Url:     urlStr,
		Params:  make(req.QueryParam),
		Headers: make(req.Header),
		expire:  time.Now().Add(DefaultSignExpireTime * time.Second),
	}
	return req, nil
}

func (req *request) buildSignKey(accessKeySecret string) {
	now := time.Now().Unix()
	expire := req.expire.Unix()
	req.SignTime = fmt.Sprintf("%d;%d", now, expire)
	req.KeyTime = req.SignTime
	req.SignKey = CreateSignature(req.KeyTime, accessKeySecret)
}

func canonicalHeaders(h http.Header) string {
	i, a, lowerCase := 0, make([]string, len(h)), make(map[string][]string)

	for k, v := range h {
		lowerCase[strings.ToLower(k)] = v
	}

	var keys []string
	for k := range lowerCase {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := lowerCase[k]
		for j, w := range v {
			v[j] = url.QueryEscape(strings.Trim(w, " "))
		}
		sort.Strings(v)
		a[i] = strings.ToLower(k) + "=" + strings.Join(v, ",")
		i++
	}
	return strings.Join(a, "&")
}

func canonicalQueryString(u *url.URL) string {
	keyValues := make(map[string]string, len(u.Query()))
	keys := make([]string, len(u.Query()))

	key_i := 0
	for k, vs := range u.Query() {
		k = strings.ToLower(k)
		k = url.QueryEscape(k)
		k = strings.ToLower(k)

		a := make([]string, len(vs))
		for idx, v := range vs {
			v = url.QueryEscape(v)
			a[idx] = fmt.Sprintf("%s=%s", k, v)
		}

		keyValues[k] = strings.Join(a, "&")
		keys[key_i] = k
		key_i++
	}

	sort.Strings(keys)

	query := make([]string, len(keys))
	for idx, key := range keys {
		query[idx] = keyValues[key]
	}

	query_str := strings.Join(query, "&")

	return strings.Replace(query_str, "+", "%20", -1)
}

func (req *request) buildFormatString() {
	FormatMethod := strings.ToLower(req.Method)

	u, err := url.Parse(req.Url)
	if err != nil {
		panic(err) // TODO log error, not panic
	}

	FormatURI := GetURIPath(u)

	urlValues := make(url.Values)
	for k, v := range req.Params {
		urlValues.Set(k, fmt.Sprint(v))
	}
	FormatParameters := urlValues.Encode()

	var headers = []string{}
	headerQuery := url.Values{}
	reqHeaders := map[string]string{
		"Host": u.Host,
	}
	for k, v := range reqHeaders {
		lowerCaseKey := strings.ToLower(k)
		headers = append(headers, lowerCaseKey)
		headerQuery.Add(lowerCaseKey, v)
	}

	sort.Strings(headers)
	req.HeaderList = strings.Join(headers, ";")

	FormatHeaders := headerQuery.Encode()

	req.FormatString = strings.Join([]string{
		FormatMethod,
		FormatURI,
		FormatParameters,
		FormatHeaders,
	}, "\n") + "\n"
}

func (req *request) buildStringToSign() {
	shaStr := MakeSha1(req.FormatString)
	req.StringToSign = strings.Join([]string{
		SignAlgorithm,
		req.SignTime,
		shaStr,
	}, "\n") + "\n"
}

func (req *request) buildSignature() {
	req.Signature = CreateSignature(req.StringToSign, req.SignKey)
}

func (req *request) buildParamList() {
	var paraKeys = []string{}
	for k, _ := range req.Params {
		lowerCaseKey := strings.ToLower(k)
		paraKeys = append(paraKeys, lowerCaseKey)
	}
	sort.Strings(paraKeys)

	req.UrlParamList = strings.Join(paraKeys, ";")
}

const (
	QSignAlgorithm = "q-sign-algorithm"
	QAK            = "q-ak"
	QSignTime      = "q-sign-time"
	QKeyTime       = "q-key-time"
	QHeaderList    = "q-header-list"
	QUrlParamList  = "q-url-param-list"
	QSign          = "q-signature"

	SignAlgorithm = "sha1"

	URLSignPara = "sign"
)

func (req *request) buildAuthorization(secretId string) {
	req.Authorization = strings.Join([]string{
		QSignAlgorithm + "=" + SignAlgorithm,
		QAK + "=" + secretId,
		QSignTime + "=" + req.SignTime,
		QKeyTime + "=" + req.KeyTime,
		QHeaderList + "=" + req.HeaderList,
		QUrlParamList + "=" + req.UrlParamList,
		QSign + "=" + req.Signature,
	}, "&")
}

func (req *request) BuildAuth(secretId, secretKey string) {
	req.buildSignKey(secretKey)
	req.buildFormatString()
	req.buildStringToSign()
	req.buildSignature()
	req.buildParamList()
	req.buildAuthorization(secretId)
}
