package main

import (
	"log"
	"os"
	"os/signal"
	"server/src"
	"server/src/components"
	"server/src/handlers"
	"sync"
	"syscall"
)

func handleExit(callback func()) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	should_exit := make(chan bool, 1)

	go func() {
		exitChan := <-sig
		log.Printf("\n######### Signal %s received. Aborting now...\n", exitChan)
		callback()
		should_exit <- true
	}()

	<-should_exit
}

func main() {
	lb := components.NewLoadBalancer(src.CONFIG_FILE_PATH).EnableLoadBalancing()

	defer handleExit(func() {
		lb.SyncFile()
	})

	midreg := components.NewMiddlewareRegistry()

	midreg.AddHandler("/", handlers.RequestService)

	var wg sync.WaitGroup
	midreg.Exec(lb, &wg)
	wg.Wait()
}
