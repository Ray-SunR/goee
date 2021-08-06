package cache

import (
	"log"
	"net/http"
	"strings"
)

type HTTPPool struct {
	self string
	basePath string
}

const defaultBasePath = "/_geecache/"

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{self: self, basePath: defaultBasePath}
}

func (httpPool *HTTPPool) Log(pattern string, values ...interface{}) {
	log.Printf(pattern, values...)
}

func (httpPool *HTTPPool) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if !strings.HasPrefix(request.URL.Path, httpPool.basePath) {
		panic("HTTPPool serving unexpected path: " + request.URL.Path)
	}

	httpPool.Log("Current URI: %s, URL.Path: %s", request.RequestURI, request.URL.Path)

	parts := strings.SplitN(request.URL.Path[len(httpPool.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(responseWriter, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(responseWriter, "no such group", http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/octet-stream")
	responseWriter.Write(view.ByteSlice())
}
