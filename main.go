package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var flagwait int

func init() {
	flag.IntVar(&flagwait, "wait", 5, "wait time in seconds for http req")
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("got request processing for", flagwait, "sec", r.RemoteAddr, r.URL)
	time.Sleep(time.Duration(flagwait) * time.Second)
	fmt.Fprintf(w, "Ok, Done after 5 sec")
}

func catchSignal(sigc chan os.Signal) {
	s := <-sigc
	log.Println("Got signal ", s.String(), "Ignoring...")

}

func main() {

	flag.Parse()

	// Don't let SIGTERM kill me
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGTERM,
	)
	go catchSignal(sigc)

	// HTTP Server
	http.HandleFunc("/", handler)
	log.Printf("WaitHTTPServer started on :8080 with %dsec wait ", flagwait)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
