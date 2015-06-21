package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
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
	proxyNames := strings.Split(strings.TrimSpace(string(dat)), ",")

	containerNameToIp := getContainersWithAddresses(proxyNames)

	// start a web server which does the API server
	fmt.Printf("going to start server with proxies %s and port %d\n", containerNameToIp, *serverPort)

	http.HandleFunc("/proxies/index", func(w http.ResponseWriter, r *http.Request) {

		jsonBytes, e := json.Marshal(proxyNames)
		if e != nil {
			http.Error(w, "JSON marshaling error", 500)
		} else {
			fmt.Fprintf(w, string(jsonBytes))
		}
	})

	http.HandleFunc("/proxies", makeProxyHandler(containerNameToIp))
	err = http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), nil)

	if err != nil {
		fmt.Println("Server start failed: ", err)
	}
}
