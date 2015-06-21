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

func corsSetup(w http.ResponseWriter, r *http.Request) bool { // CORS header
	origin := r.Header["Origin"]

	if len(origin) == 1 {
		w.Header().Add("Access-Control-Allow-Origin", origin[0])
	}

	w.Header().Add("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,"+
		"If-Modified-Since,Cache-Control,Content-Type,Referer,x-access-token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)

	if r.Method != "OPTION" {
		return false
	}
	return true
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

	proxyNamesRaw := strings.Split(strings.TrimSpace(string(dat)), "\n")
	proxyNames := make([]string, len(proxyNamesRaw))
	copy(proxyNames, proxyNamesRaw)
	for idx, proxyName := range proxyNamesRaw {
		parts := strings.Split(proxyName, ":")
		if len(parts) == 2 {
			fmt.Println("IDX: ", idx)
			proxyNames[idx] = "figtest_" + parts[0] + parts[1] + "_1"
		} else {
			proxyNames[idx] = parts[0] + "_" + parts[1] + parts[2] + "_1"
		}
	}

	containerNameToIp := getContainersWithAddresses(proxyNames)

	// start a web server which does the API server
	fmt.Printf("going to start server with proxies %s and port %d\n", containerNameToIp, *serverPort)

	http.HandleFunc("/proxies/index", func(w http.ResponseWriter, r *http.Request) {
		if corsSetup(w, r) {
			return
		}

		// Return the raw file to allow frontend to understand topology.
		jsonBytes, e := json.Marshal(proxyNamesRaw)

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
