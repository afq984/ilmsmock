package proxy

import (
	"io"
	"log"
	"net/http"
)

const (
	metaName = "meta.json"
	bodyName = "body"
)

func errorText(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

type ServerOptions struct {
	UpstreamPrefix string
	CaptureRoot    string
}

type Server struct {
	Upstream
}

var _ http.Handler = &Server{}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath, cacheKey := Key(r.URL)
	cachedResponse, err := s.Upstream.Get(requestPath, cacheKey)
	if err != nil {
		log.Println(err)
		errorText(w, http.StatusNotFound)
		return
	}

	CopyHeader(w.Header(), cachedResponse.Headers())
	w.WriteHeader(cachedResponse.Status())

	io.Copy(w, cachedResponse)
	cachedResponse.Close()
}
