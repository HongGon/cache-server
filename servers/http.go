package servers

import (
	"cache-server/caches"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HTTPServer is HTTP server structure
type HTTPServer struct {
	// Cache is underlying structure to storage
	cache *caches.Cache
}

// NewHTTPServer returns a new HTTPServer
func NewHTTPServer(cache *caches.Cache) *HTTPServer {
	return &HTTPServer{
		cache: cache,
	}
}

// launch server
func (hs *HTTPServer) Run(address string) error {
	return http.ListenAndServe(address, hs.routerHandler())
}


// routerHandler return routerHandler
func (hs *HTTPServer) routerHandler() http.Handler {
	// httprouter.New() create a http router, including all request methods
	// GET request method is used to query the cache
	// PUT request method is used to create a cache
	// DELETE request method is used to del the cache

	// key is obtained from url, value is obtained from request 
	router := httprouter.New()
	router.GET("/cache/:key", hs.getHandler)
	router.PUT("/cache/:key", hs.setHandler)
	router.DELETE("/cache/:key", hs.deleteHandler)
	router.GET("/status", hs.statusHandler)
	return router
}

// getHandler is used to obtain cache data
func (hs *HTTPServer) getHandler(writer http.ResponseWriter,
	request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	value, ok := hs.cache.Get(key)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.Write(value)
}


// setHandler is used to storage cache data
func (hs *HTTPServer) setHandler(writer http.ResponseWriter, 
	request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	value, err := ioutil.ReadAll(request.Body)
	if err != nil {
		//  if fail, return 500
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	hs.cache.Set(key, value)
}

// deleteHandler is used to delete cache data
func (hs *HTTPServer) deleteHandler(writer http.ResponseWriter, 
	request *http.Request, params httprouter.Params) {
	key := params.ByName("key")
	hs.cache.Delete(key)
}

// statusHandler is used to obtain the num of k-v
func (hs *HTTPServer) statusHandler(writer http.ResponseWriter,
	request *http.Request, params httprouter.Params) {
	// encoding num to JSON string
	status, err := json.Marshal(map[string]interface{}{
		"count": hs.cache.Count(),
	})
	if err != nil {
		// if encoding fails, return 500 status code
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(status)
}




















