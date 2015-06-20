package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type proxyEndpoint struct {
	containerName string
	port          int
}

func makeProxyEndpoint(proxyString string) proxyEndpoint {
	proxyComponents := strings.Split(proxyString, ":")
	port, _ := strconv.Atoi(proxyComponents[1])
	return proxyEndpoint{proxyComponents[0], port}
}

func main() {
	filename := flag.String("filename", "", "file containing list of proxies")
	serverPort := flag.Int("port", 8080, "port for API server")
	flag.Parse()

	if *filename == "" {
		fmt.Println("filename must be specified")
		os.Exit(1)
	}

	dat, err := ioutil.ReadFile(*filename)
	check(err)
	proxyStrings := strings.Split(strings.TrimSpace(string(dat)), ",")

	var proxies = make([]proxyEndpoint, len(proxyStrings))
	for i, proxyString := range proxyStrings {
		proxy := makeProxyEndpoint(proxyString)
		proxies[i] = proxy
	}

	// start a web server which does the API server
	fmt.Printf("going to start server with proxies %s and port %d\n", proxies, *serverPort)

	http.HandleFunc("/proxies/index", func(w http.ResponseWriter, r *http.Request) {
		var proxyList = make([]string, len(proxies))
		for i, proxy := range proxies {
			proxyList[i] = proxy.containerName
		}

		jsonBytes, e := json.Marshal(proxyList)
		if e != nil {
			http.Error(w, "JSON marshaling error", 500)
		} else {
			fmt.Fprintf(w, string(jsonBytes))
		}
	})

	http.HandleFunc("/proxies", makeProxyHandler())
	err = http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), nil)

	if err != nil {
		fmt.Println("Server start failed: ", err)
	}
}
