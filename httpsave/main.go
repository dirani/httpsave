package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	var listenAddr string
	var saveFile string

	flag.StringVar(&listenAddr, "listenAddr", ":8080", "listen address")
	flag.StringVar(&saveFile, "saveFile", "/tmp/httpsave", "saveFile")
	flag.Parse()

	log.Printf("installing default catch all handler at /")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { rootHandler(w, r) })

	log.Printf("installing save handler at /save")
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) { saveHandler(w, r, saveFile) })
	http.HandleFunc("/save/", func(w http.ResponseWriter, r *http.Request) { saveHandler(w, r, saveFile) })

	log.Printf("server will listen at %s", listenAddr)
	server := &http.Server{Addr: listenAddr}

	log.Fatal(server.ListenAndServe())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	const notFound = "404 nothing here\n"
	log.Printf("rootHandler: method=%s host=%s path=%s from=%s - response:%s", r.Method, r.Host, r.URL.Path, r.RemoteAddr, notFound)
	http.Error(w, notFound, 404)
}

func saveHandler(w http.ResponseWriter, r *http.Request, saveFile string) {

	const internalError = "500 internal server error\n"

	body, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		log.Printf("saveHandler: method=%s host=%s path=%s from=%s - error:%s response:%s", r.Method, r.Host, r.URL.Path, r.RemoteAddr, errRead, internalError)
		http.Error(w, internalError, 500)
		return
	}

	log.Printf("saveHandler: method=%s host=%s path=%s from=%s - saving %d bytes to file %s", r.Method, r.Host, r.URL.Path, r.RemoteAddr, len(body), saveFile)

	errWrite := ioutil.WriteFile(saveFile, body, 0640)
	if errWrite != nil {
		log.Printf("saveHandler: method=%s host=%s path=%s from=%s - error:%s response:%s", r.Method, r.Host, r.URL.Path, r.RemoteAddr, errWrite, internalError)
		http.Error(w, internalError, 500)
		return
	}

	const responseOk = "200 ok\n"
	log.Printf("saveHandler: method=%s host=%s path=%s from=%s - response:%s", r.Method, r.Host, r.URL.Path, r.RemoteAddr, responseOk)
	io.WriteString(w, responseOk)
}
