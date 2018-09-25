package cos

import (
	"fmt"
)

type Option func(*request) error

func Prefix(prefix string) Option {
	return addParam("prefix", prefix)
}

func Delimiter(delimiter string) Option {
	return addParam("delimiter", delimiter)
}

func Marker(marker string) Option {
	return addParam("marker", marker)
}

func MaxKeys(maxKeys int) Option {
	return addParam("max-keys", maxKeys)
}

func addParam(key string, value interface{}) Option {
	return func(req *request) error {
		if value == nil {
			return nil
		}
		req.Params[key] = value
		return nil
	}
}

func ContentType(contentType string) Option {
	return addHeader("Content-Type", contentType)
}
func addHeader(key string, value interface{}) Option {
	return func(req *request) error {
		if value == nil {
			return nil
		}
		req.Headers[key] = fmt.Sprint(value)
		return nil
	}
}

func Body(body interface{}) Option {
	return func(req *request) error {
		req.Payload = body
		return nil
	}
}
