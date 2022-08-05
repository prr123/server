package main

import (
	"fmt"
	"net/http"
	"os"
)

type httpObj struct {
	dbg bool
}

func initHttp()(sObj *httpObj) {

	sObj = new(httpObj)
	sObj.dbg = false

	return sObj
}

func (sObj *httpObj) handle(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("dbg: %t\n", sObj.dbg)
	// Log the request protocol
	if sObj.dbg {
	fmt.Printf("\n *** Got connection: %s from: %v\n",r.Proto, r.RemoteAddr)
	fmt.Println("  Request URI: ", r.RequestURI)
	fmt.Println("  URL: ", r.URL, " Host: ", r.Host)
	fmt.Println("  Method   : ", r.Method)

	fmt.Println(" *** headers *** ")
	for k, v := range r.Header {
		fmt.Printf("  %-25s : %-25s \n", k, v)
	}
	}

// Send a message back to the client
	//	w.Write([]byte("Hello"))
	p := "./static/index.html"
    http.ServeFile(w, r, p)
}


func main() {

	argsNum := len(os.Args)

	// initialize server  object
	sObj := initHttp()
	if sObj == nil {
		fmt.Printf("error init sObj!\n")
		os.Exit(-1)
	}

	switch argsNum {
		case 1:

		case 2:
			if os.Args[1] == "dbg" {sObj.dbg = true}

		default:
			fmt.Println("error cmd line too many args!")
			os.Exit(-1)
	}


	// Create a server on port 8000
	// Exactly how you would run an HTTP/1.1 server
	srv := &http.Server{Addr: "127.0.0.1:8002", Handler: http.HandlerFunc(sObj.handle)}

	// Start the server with TLS, since we are running HTTP/2 it must be
	// run with TLS.
	// Exactly how you would run an HTTP/1.1 server with TLS connection.
	fmt.Printf("Serving localhost on https: 127.0.0.1:8002 \n")
	err := srv.ListenAndServeTLS("/home/peter/newca/server_crt.pem", "/home/peter/newca/server_key.pem")
	if err != nil {
		fmt.Printf("error starting server: %v\n", err)
		os.Exit(-11)
	}
}

