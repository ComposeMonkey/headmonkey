package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const BEHAVIOR_URL = "http://%s/behavior"
const BEHAVIOR_BODY = `{ "behavior": "%s"}`

type handler func(http.ResponseWriter, *http.Request)

type VConfig struct {
	Behavior string `json:"behavior"`
}

func updateVConfig(proxyName string, behavior string) error {
	jsonStr := []byte(fmt.Sprintf(BEHAVIOR_BODY, behavior))
	url := fmt.Sprintf(proxyName, BEHAVIOR_URL)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	return err
}

func getVConfig(proxyName string) (error, *VConfig) {
	resp, err := http.Get(fmt.Sprintf(proxyName, BEHAVIOR_URL))
	defer resp.Body.Close()

	if err != nil {
		return err, nil
	}

	vconfig := VConfig{}

	rawJSON, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rawJSON, &vconfig)

	return nil, &vconfig
}

func getProxyName(uri string) string {
	urlParts := strings.Split(uri, "/")
	proxyName := urlParts[len(urlParts)-1]
	return proxyName
}

func makeGETProxyHandler() handler {
	return func(res http.ResponseWriter, req *http.Request) {
		proxyName := getProxyName(req.URL.RequestURI())
		if err, vconfig := getVConfig(proxyName); err == nil {
			reply, _ := json.Marshal(*vconfig)
			res.Write(reply)
		}
	}
}

func makePUTProxyHandler() handler {
	return func(res http.ResponseWriter, req *http.Request) {
		proxyName := getProxyName(req.URL.RequestURI())
		rawJSON, _ := ioutil.ReadAll(req.Body)
		err := updateVConfig(proxyName, string(rawJSON))

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
			res.Write([]byte("BAD"))
		}
		res.Write([]byte("OK"))
	}
}
