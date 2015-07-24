package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vallard/twitter-client/twitstream"
)

func main() {
	a := os.Args[1:]

	if len(a) < 1 {
		fmt.Println("first argument should be a search value")
		return
	}

	// the stream we want to listen for with the search terms.
	streamMe := strings.Join(a, ",")
	fmt.Println(streamMe)

	tc := twitstream.New()
	stop := false
	signalChan := make(chan os.Signal, 1)

	go func() {
		<-signalChan
		stop = true
		log.Println("Stopping...")
		tc.CloseConn()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		if stop {
			return
		}
		tc.Get(streamMe)
		// sleep to avoid getting twitter upset with us.
		time.Sleep(2 * time.Second)
	}

}
