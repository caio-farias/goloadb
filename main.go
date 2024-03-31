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
		log.Println("Waiting exit signal.")
		exitChan := <-sig
		log.Printf("\n######### Signal %s received. Aborting now...\n", exitChan)
		callback()
		should_exit <- true
	}()

	<-should_exit
}

func main() {
	lb := components.NewLoadBalancer(src.CONFIG_FILE_PATH)
	lb.EnableLoadBalancing()

	midreg := components.NewMiddlewareRegistry()
	midreg.AddHandler("/", handlers.RequestService)

	defer handleExit(func() {
		lb.SyncFile()
	})

	var wg sync.WaitGroup
	wg.Add(1)
	midreg.Exec(lb, &wg)
	wg.Wait()
}
