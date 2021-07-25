package proxy

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"path"
)

type meta struct {
	Status  int
	Headers http.Header
}

func newMeta() meta {
	return meta{
		Headers: make(http.Header),
	}
}

func Key(in *url.URL) (requestPath, cacheKey string) {
	out := &url.URL{
		Path:     path.Clean(in.Path),
		RawQuery: in.Query().Encode(),
	}
	requestPath = out.String()
	cacheKey = base64.RawURLEncoding.EncodeToString(
		[]byte(requestPath),
	)
	return
}

type cachedResponse interface {
	io.ReadCloser
	Headers() http.Header
	Status() int
}

type cachedResponseImpl struct {
	io.ReadCloser
	meta
}

var _ cachedResponse = &cachedResponseImpl{}

func (r *cachedResponseImpl) Headers() http.Header {
	return r.meta.Headers
}

func (r *cachedResponseImpl) Status() int {
	return r.meta.Status
}

var importantHeaders = []string{
	"Content-Type",
	"Content-Disposition",
	"Location",
}

func CopyHeader(dst, src http.Header) {
	for _, name := range importantHeaders {
		value := src.Get(name)
		if value != "" {
			dst.Set(name, value)
		}
	}
}
