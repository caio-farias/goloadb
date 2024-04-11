package main

import (
	"net/http"
	"server/src/components"
)

func main() {
	lb := components.NewLoadBalancer("./config.json")

	sv := components.NewServer("3000")

	sv.EnableLoadBalancing(lb, "/")

	sv.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte(sv.GetDiscoveryInfo()))
		return nil
	})

	sv.Run()
}
