package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const BEHAVIOR_URL = "http://%s/behavior"

type handler func(http.ResponseWriter, *http.Request)

type VConfig struct {
	Behavior string `json:"behavior"`
}

func updateVConfig(proxyName string, behaviorJSON string) error {
	jsonStr := []byte(behaviorJSON)
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

func getProxyName(req *http.Request) string {
	proxyName := req.URL.Query().Get("proxy")
	return proxyName
}

func GETProxyHandler(res http.ResponseWriter, req *http.Request) {
	proxyName := getProxyName(req)
	if err, vconfig := getVConfig(proxyName); err == nil {
		reply, _ := json.Marshal(*vconfig)
		res.Write(reply)
	}
}

func PUTProxyHandler(res http.ResponseWriter, req *http.Request) {
	proxyName := getProxyName(req)
	rawJSON, _ := ioutil.ReadAll(req.Body)
	err := updateVConfig(proxyName, string(rawJSON))

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		res.Write([]byte("BAD"))
	}
	res.Write([]byte("OK"))
}

func makeProxyHandler() handler {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "PUT" {
			PUTProxyHandler(res, req)
		} else if req.Method == "GET" {
			GETProxyHandler(res, req)
		}
	}
}
