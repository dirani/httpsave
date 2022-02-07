package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	xj "github.com/basgys/goxml2json"
)

func main() {

	var listenAddr string
	var saveFile string
	var saveJson string

	flag.StringVar(&listenAddr, "listenAddr", ":8080", "listen address [host]:port")
	flag.StringVar(&saveFile, "saveFile", "C:\\temp\\httpsave", "save file path")
	flag.StringVar(&saveJson, "saveJson", "C:\\temp\\httpsave.json", "save json file path")
	flag.Parse()

	log.Printf("installing default catch all handler at /")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { rootHandler(w, r) })

	log.Printf("installing save handler at /save")
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) { saveHandler(w, r, saveFile) })
	http.HandleFunc("/save/", func(w http.ResponseWriter, r *http.Request) { saveHandler(w, r, saveFile) })
	http.HandleFunc("/x2j", func(w http.ResponseWriter, r *http.Request) { x2jHandler(w, r, saveJson) })
	http.HandleFunc("/x2j/", func(w http.ResponseWriter, r *http.Request) { x2jHandler(w, r, saveJson) })

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

func x2jHandler(w http.ResponseWriter, r *http.Request, saveFile string) {

	const me = "x2jHandler"

	const internalError = "500 internal server error\n"

	body, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		log.Printf("%s: method=%s host=%s path=%s from=%s - error:%s response:%s", me, r.Method, r.Host, r.URL.Path, r.RemoteAddr, errRead, internalError)
		http.Error(w, internalError, 500)
		return
	}

	log.Printf("%s: method=%s host=%s path=%s from=%s - received %d bytes", me, r.Method, r.Host, r.URL.Path, r.RemoteAddr, len(body))

	xml := strings.NewReader(string(body))
	j, errConv := xj.Convert(xml)
	if errConv != nil {
		log.Printf("%s: method=%s host=%s path=%s from=%s - error:%s response:%s", me, r.Method, r.Host, r.URL.Path, r.RemoteAddr, errConv, internalError)
		http.Error(w, internalError, 500)
		return
	}

	jBytes := j.Bytes()

	log.Printf("%s: method=%s host=%s path=%s from=%s - saving JSON %d bytes to file %s", me, r.Method, r.Host, r.URL.Path, r.RemoteAddr, len(jBytes), saveFile)

	errWrite := ioutil.WriteFile(saveFile, jBytes, 0640)
	if errWrite != nil {
		log.Printf("%s: method=%s host=%s path=%s from=%s - error:%s response:%s", me, r.Method, r.Host, r.URL.Path, r.RemoteAddr, errWrite, internalError)
		http.Error(w, internalError, 500)
		return
	}

	const responseOk = "200 ok\n"
	log.Printf("%s: method=%s host=%s path=%s from=%s - response:%s", me, r.Method, r.Host, r.URL.Path, r.RemoteAddr, responseOk)
	io.WriteString(w, responseOk)
}
