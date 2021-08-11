package cache

import (
	"consistenthash"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplica  = 50
)

type HTTPPool struct {
	self        string
	basePath    string
	mutex       sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter
}

type httpGetter struct {
	baseURL string
}

func (h *HTTPPool) Set(peers ...string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.peers = consistenthash.New(defaultReplica, nil)
	h.peers.AddPhysicalNodes(peers...)
	h.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		h.httpGetters[peer] = &httpGetter{baseURL: peer + h.basePath}
	}
}

func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if peer := h.peers.GetPhysicalNode(key); peer != "" && peer != h.self {
		h.Log("Pick peer %s", peer)
		return h.httpGetters[peer], true
	}
	return nil, false
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	url := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

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
