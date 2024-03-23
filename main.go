package main

import (
	"log"
	"os"
	"os/signal"
	"server/src"
	"server/src/components"
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
	lb := components.NewLoadBalancer(src.CONFIG_FILE_PATH)

	defer handleExit(func() {
		lb.SyncFile()
	})

	midreg := components.NewMiddlewareRegistry()

	midreg.Add("/", func(url string, header *components.Header, body string) (string, error) {
		request, err := components.NewRequest(url, header)
		if err != nil {
			log.Println(err)
			return "", err
		}
		res, err := request.Await()
		if err != nil {
			log.Println(err)
			return "", err
		}

		return res, nil
	})

	var wg sync.WaitGroup
	midreg.Exec(lb, &wg)
	wg.Wait()
}
