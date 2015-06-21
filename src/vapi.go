package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const BEHAVIOR_URL = "http://%s:2020/behavior"

type handler func(http.ResponseWriter, *http.Request)

type VConfig struct {
	Behavior string `json:"behavior"`
}

func updateVConfig(proxyName string, behaviorJSON string) error {
	jsonStr := []byte(behaviorJSON)
	url := fmt.Sprintf(BEHAVIOR_URL, proxyName)
        fmt.Println("PUT: %s", behaviorJSON)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return err
}

func getVConfig(proxyName string) (error, *VConfig) {
	resp, err := http.Get(fmt.Sprintf(BEHAVIOR_URL, proxyName))
        fmt.Println("ERROR: ", err)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	vconfig := VConfig{}

	rawJSON, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rawJSON, &vconfig)

	return nil, &vconfig
}

func getProxyName(req *http.Request, containerNameToIp map[string]string) string {
	proxyName := req.URL.Query().Get("proxy")
	ipAddress := containerNameToIp["/" + proxyName]

	fmt.Printf("translating container name %s to ip address %s using map %s\n", proxyName, ipAddress, containerNameToIp)
	return ipAddress
}

func GETProxyHandler(res http.ResponseWriter, req *http.Request, containerNameToIp map[string]string) {
	proxyName := getProxyName(req, containerNameToIp)
	fmt.Printf("sending GET to %s\n", proxyName)
	if err, vconfig := getVConfig(proxyName); err == nil {
		reply, _ := json.Marshal(*vconfig)
		res.Write(reply)
	}
}

func PUTProxyHandler(res http.ResponseWriter, req *http.Request, containerNameToIp map[string]string) {
	proxyName := getProxyName(req, containerNameToIp)
	rawJSON, _ := ioutil.ReadAll(req.Body)
	err := updateVConfig(proxyName, string(rawJSON))

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		res.Write([]byte("BAD"))
	}
	res.Write([]byte("OK"))
}

func makeProxyHandler(containerNameToIp map[string]string) handler {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "PUT" {
			PUTProxyHandler(res, req, containerNameToIp)
		} else if req.Method == "GET" {
			GETProxyHandler(res, req, containerNameToIp)
		}
	}
}
