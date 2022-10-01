package main

import (
	"cache-server/caches"
	"cache-server/servers"
)

func main() {
	caches := caches.NewCache()
	err := servers.NewHTTPServer(caches).Run(":5837")
	if err != nil {
		panic(err)
	}
}







