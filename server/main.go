package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"

	_ "net/http/pprof"
)

func main() {
	os.Open("xyz")
	var listenPort = flag.String("p", ":8080", "port to listen")
	flag.Parse()
	// start pprof
	runtime.SetBlockProfileRate(1)
	go func() {
		log.Println("http://localhost:6060/debug/pporf")
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	// server main
	listner, err := net.Listen("tcp", *listenPort)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			defer conn.Close()
			log.Printf("Accept %v\n", conn.RemoteAddr())
			request, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				log.Fatal(err)
			}
			dump, err := httputil.DumpRequest(request, true)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(dump))
			response := http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Body:       ioutil.NopCloser(strings.NewReader("Hello World\n")),
			}
			response.Write(conn)
		}()
	}
}
