package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Upstream interface {
	Get(requestPath, cacheKey string) (cachedResponse, error)
}

type FileSystemUpstream struct {
	fs fs.FS
}

var _ Upstream = &FileSystemUpstream{}

func NewFileSystemUpstream(fs fs.FS) *FileSystemUpstream {
	return &FileSystemUpstream{
		fs: fs,
	}
}

func (up *FileSystemUpstream) Get(requestPath, cacheKey string) (cachedResponse, error) {
	jsonb, err := fs.ReadFile(up.fs, filepath.Join(cacheKey, metaName))
	if err != nil {
		return nil, err
	}

	meta := newMeta()
	err = json.Unmarshal(jsonb, &meta)
	if err != nil {
		return nil, err
	}

	body, err := up.fs.Open(filepath.Join(cacheKey, bodyName))
	if err != nil {
		return nil, err
	}
	return &cachedResponseImpl{
		ReadCloser: body,
		meta:       meta,
	}, nil
}

type HTTPUpstream struct {
	origin      string
	captureRoot string
	client      http.Client
}

var _ Upstream = &HTTPUpstream{}

func NewHTTPUpstream(origin, cacheRoot string) *HTTPUpstream {
	return &HTTPUpstream{
		origin:      origin,
		captureRoot: cacheRoot,
		client: http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (s *HTTPUpstream) Get(requestPath, cacheKey string) (cachedResponse, error) {
	upstreamResp, err := s.client.Get(s.origin + requestPath)
	if err != nil {
		log.Panic(err)
	}
	if upstreamResp.StatusCode >= 500 {
		log.Panicf("Unexcepted status %d %v", upstreamResp.StatusCode, requestPath)
	}

	meta := newMeta()
	meta.Status = upstreamResp.StatusCode
	CopyHeader(meta.Headers, upstreamResp.Header)

	var b bytes.Buffer
	_, err = io.Copy(&b, upstreamResp.Body)
	if err != nil {
		panic(err)
	}

	cacheDir := filepath.Join(s.captureRoot, cacheKey)
	os.MkdirAll(cacheDir, 0755)

	jsonb, err := json.Marshal(&meta)
	if err != nil {
		log.Panic(err)
	}
	err = os.WriteFile(filepath.Join(cacheDir, metaName), jsonb, 0644)
	if err != nil {
		log.Panic(err)
	}
	err = os.WriteFile(filepath.Join(cacheDir, bodyName), b.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}

	return &cachedResponseImpl{
		ReadCloser: io.NopCloser(bytes.NewReader(b.Bytes())),
		meta:       meta,
	}, nil
}
