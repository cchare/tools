package main

// a simple file server
import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func fserver(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Listen on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir("."))))
}

func server(host string, port int, content string) {
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Listen on %s\n", addr)
	body := []byte(content)

	http.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	})

	http.HandleFunc("/a/b", func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}

func router(host string, port int, content string) {
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Listen on %s\n", addr)
	body := []byte(content)
	router := httprouter.New()

	router.Handle("GET", "/a", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		for k, v := range r.Header {
			fmt.Printf("%s : %v\n", k, v)
		}

		htest := r.Header.Get("x-htest")
		fmt.Printf("header-> x-htest[%d]: %s\n", len(htest), htest)

		ht, ok := r.Header["X-Htest"]
		if ok {
			fmt.Printf("ok : %v\n", ht)
		} else {
			fmt.Printf("not: %v\n", ht)
		}

		w.Write(body)
	})

	router.Handle("GET", "/a/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write(body)
	})

	log.Fatal(http.ListenAndServe(addr, router))
}

func main() {
	host := flag.String("host", "", "Listen address")
	port := flag.Int("port", 80, "Listen port")
	data := flag.String("data", "default http body", "http response body")

	flag.Parse()

	//fserver(*host, *port)

	//server(*host, *port, *data)

	m := map[string]string{}
	m["0"] = ""
	m["1"] = "a"
	fmt.Printf("%#v\n", m)

	router(*host, *port, *data)
}
