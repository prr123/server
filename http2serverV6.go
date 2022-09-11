// http2serverV5.go
// program that creates the directory structure for a domain
//
// author: prr, azul software
// date: 5 Septtember 2022
// copyright 2022 prr, azul software
//
// ver 5 add js css routes
//

package main

import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
)

type httpObj struct {
	dbg bool
	baseUri string
	idxPath string
}

type tstform struct {
	Formtyp string `json:"form"`
	Formdat formdata `json:"data"`
}

type formdata struct {
	Usrnam string `json:"usrnam"`
	Usremail string `json:"usremail"`
}

func initHttp()(sObj *httpObj) {

	sObj = new(httpObj)
	sObj.dbg = false
	return sObj
}

func getFileNames() (jData *[]byte, err error) {

    target := "/home/peter/www/azultest/html"

    dir, err := os.Open(target)
    if err != nil { return nil, fmt.Errorf("cannot open target dir! %v\n", err)}
    defer dir.Close()

    // verify directory


    // get files

    files, err := dir.ReadDir(0)
    if err != nil {return nil, fmt.Errorf("cannot read files in dir! %v\n", err)}

    count:=0
	fileCount := 10 //len(files)
	fileArr := make([]string, fileCount)
    for i:=0; i< fileCount; i++ {
        if files[i].IsDir() {continue}
        filnam := target + "/" + files[i].Name()
        info, err := os.Stat(filnam)
        if err != nil {return nil, fmt.Errorf("error getting filestat for %s! %v\n", files[i].Name(), err)}
        t := info.ModTime()
//      timeStr:= t.Format(time.RFC822)
        timeStr:= t.Format("02/01/2006 15:04:05")
        fmt.Printf("file %d: %-20s %4dK %s\n",i, files[i].Name(), info.Size()/1000, timeStr)
		fileArr[i] = files[i].Name()
        count++
    }
    json, err := json.Marshal(fileArr)
    if err != nil {return nil, fmt.Errorf("error marshal: %v", err)}
	return &json, nil
}


func (sObj *httpObj) handle(w http.ResponseWriter, r *http.Request) {

	// Log the request protocol
	if sObj.dbg {
		fmt.Printf(" *** Received connection request from: %v from: %s\n",r.RemoteAddr, r.Proto)
		fmt.Println("  Request URI:    ", r.RequestURI)
		fmt.Println("  Request URL:    ", r.URL, " Host: ", r.Host)
		fmt.Println("  Request Method: ", r.Method)

		fmt.Println("  *** headers *** ")
		for k, v := range r.Header {
			fmt.Printf("    %-25s : %-25s \n", k, v)
		}
	}

	// Send a message back to the client
	//	w.Write([]byte("Hello"))
	// now I can pass parameters via sObj

	if r.Method == "GET" {
		p := r.RequestURI
		if p == "/" {
			p = sObj.idxPath
		}

		extPos:=-1
		for i:=len(p)-1; i>=0; i-- {
			if p[i] == '.' {
				extPos = i
			}
		}

		if extPos == -1 {
			p += ".html"
			extPos = len(p)
		}

		sFil:=""
//		fileJs:=false
		ext := string(p[extPos+1:])
		fmt.Printf("ext: %s\n", ext)
		switch ext {
			case "html":
				sFil = "/html" + p
			case "pdf":
				sFil = "/pdf" + p
			case "css":
				sFil = "/css" + p
			case "js":
				sFil = "/js" + p

			case "json":
//				fileJs=true
				jData, err := getFileNames()
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
//					w.WriteHeader(http.StatusNotAcceptable) // 406
 					http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated) // added by prr
				w.Write(*jData)
				return
			default:
				fmt.Printf("invalid extension: %s\n", ext)
				return
		}



		if p[:6] != "/home/" {
			p= sObj.baseUri + sFil
		}

		fmt.Printf("\n*** serving %s ***\n",p)
    	http.ServeFile(w, r, p)
		return
	}

	if r.Method == "POST" {
		p := r.RequestURI
		ct := r.Header.Get("Content-Type")
		fmt.Printf("\n*** post received: %s content: %s ***\n", p, ct)

//		r.ParseForm()
//		if len(r.Form)

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("error reading reqbody: %v", err)
		}

		fmt.Printf("body:\n%s\n", reqBody)
		var replForm tstform
//		fmt.Printf("testForm:\n%v\n", tstForm)
		err = json.Unmarshal(reqBody, &replForm)
		if err != nil {
			fmt.Printf("error - unmarshal: %v\n", err)
		}
		fmt.Printf("reply form:/n%+v\n", replForm)
		return
	}

}


func main() {

	var serveUri string
	argsNum := len(os.Args)

	// initialize server  object
	sObj := initHttp()
	if sObj == nil {
		fmt.Printf("error init sObj!\n")
		os.Exit(-1)
	}
	baseUri :="/home/peter/www/"

	switch argsNum {
		case 1:
			serveUri = baseUri + "base"
		case 2:
			if os.Args[1] == "dbg" {
				sObj.dbg = true
			} else {
				serveUri = baseUri + os.Args[1]
			}
		case 3:
			serveUri = baseUri + os.Args[1]
			if os.Args[2] == "dbg" {sObj.dbg = true}

		default:
			fmt.Println("error cmd line too many args!")
			os.Exit(-1)
	}

	_,err := os.Stat(serveUri)
	if err != nil {
		fmt.Printf("error %v\n", err)
		os.Exit(-1)
	}


	sObj.baseUri = serveUri
	sObj.idxPath = serveUri + "/html/index.html"

	fmt.Printf("Uri is %s\nIndex File is %s\n", sObj.baseUri, sObj.idxPath)

	// Create a server on port 8000
	// Exactly how you would run an HTTP/1.1 server
	srv := &http.Server{Addr: "127.0.0.1:8002", Handler: http.HandlerFunc(sObj.handle)}

	// Start the server with TLS, since we are running HTTP/2 it must be
	// run with TLS.
	// Exactly how you would run an HTTP/1.1 server with TLS connection.
	fmt.Printf("Serving localhost on https: 127.0.0.1:8002 \n")
	err = srv.ListenAndServeTLS("/home/peter/newca/server_crt.pem", "/home/peter/newca/server_key.pem")
	if err != nil {
		fmt.Printf("error starting server: %v\n", err)
		os.Exit(-11)
	}
}

