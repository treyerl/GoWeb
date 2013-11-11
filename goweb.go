package main

import (
	"fmt"
	"flag"
	"net/http"
  	"github.com/wsxiaoys/terminal/color"
  	"strings"
  	"strconv"
)

type fileCheckHandler struct {
	root http.FileSystem
	handler func(http.ResponseWriter, *http.Request)
	index func(http.ResponseWriter, *http.Request)
}

func FileCheckServer(root http.FileSystem, handler func(http.ResponseWriter, *http.Request), index func(http.ResponseWriter, *http.Request)) http.Handler{
	return &fileCheckHandler{root, handler, index}
}

func (f *fileCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	// sanitization : http://golang.org/src/pkg/net/http/fs.go?s=12008:12048#L401
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	
	// index check
	if r.URL.Path == "/" {
		index, err := f.root.Open("index.html")
		if err != nil {
			if f.index != nil {
				f.index(w, r)
			} else {
				fmt.Fprint(w, "This is your Website! No index file or index handler configured.")
			}
		} else {
			d, _ := index.Stat()
			http.ServeContent(w, r, d.Name(), d.ModTime(), index)
		}
		
	// file check
	} else {
		file, err := f.root.Open(upath)
		if err != nil {
			if f.handler != nil {
				f.handler(w, r)
			} else {
				http.NotFound(w, r)
			}
		} else {
			d, _ := file.Stat()
			http.ServeContent(w, r, d.Name(), d.ModTime(), file)
		}
	}
}


func test(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func start(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "This is the homepage! %s", r.URL.Path[1:])
}

func main(){
	var webroot = flag.String("webroot", "/tmp/", "The folder being served to the web.")
	var port = flag.Int("port", 8080, "The port the server is listening to.")
	flag.Parse()
	color.Printf("@yWebroot: %s, Port: %d\r\n",*webroot, *port)
	http.Handle("/",FileCheckServer(http.Dir(*webroot), nil, nil))
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}