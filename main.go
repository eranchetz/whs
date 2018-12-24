package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var flagwait int
var flagIgnoresSIGTERM bool

type request struct {
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}

func init() {
	flag.IntVar(&flagwait, "wait", 0, "wait time in seconds for http req")
	flag.BoolVar(&flagIgnoresSIGTERM, "ignore-sigterm", true, "ignore SIGTERM [default=true]")
}

// middlewareWait makes the response wait, default is 0 you can set it by requesting:
// GET /?wait=5s
func middlewareWait(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		duration := time.Duration(flagwait) * time.Second
		var err error

		u := r.URL
		qp, _ := url.ParseQuery(u.RawQuery)
		if qp["wait"] != nil {
			duration, err = time.ParseDuration(qp["wait"][0])
			if err != nil {
				fmt.Fprintf(w, "error parsing wait %v", err)
			}
		}
		log.Println("got request processing for ", duration, r.RemoteAddr, r.URL)
		time.Sleep(duration)
		next.ServeHTTP(w, r)
	})
}

// handlerHealth returns ok
func handlerHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

// handlerEcho returns the request headers and body in json form
func handlerEcho(rw http.ResponseWriter, r *http.Request) {
	var err error
	rr := &request{}
	rr.Method = r.Method
	rr.Headers = r.Header
	rr.URL = r.URL.String()
	rr.Body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rrb, err := json.Marshal(rr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(rrb)
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
	if flagIgnoresSIGTERM {
		go catchSignal(sigc)
	}
	// HTTP Server
	handleEcho := http.HandlerFunc(handlerEcho)
	handleHealth := http.HandlerFunc(handlerHealth)
	http.Handle("/", middlewareWait(handleEcho))
	http.Handle("/health", middlewareWait(handleHealth))

	log.Printf("WaitHTTPServer started on :8080 with %dsec wait , PID=%d, IgnoresSIGTERM=%t", flagwait, os.Getpid(), flagIgnoresSIGTERM)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
