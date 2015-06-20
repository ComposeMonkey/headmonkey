package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const BEHAVIOR_URL = "http://%s/behavior"
const BEHAVIOR_BODY = `{ "behavior": "%s"}`

type VConfig struct {
	Behavior string `json:"behavior"`
}

func updateVConfig(proxyUrl string, behavior string) error {
	jsonStr := []byte(fmt.Sprintf(BEHAVIOR_BODY, behavior))
	url := fmt.Sprintf(proxyUrl, BEHAVIOR_URL)
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

func getVConfig(proxyUrl string) (error, *VConfig) {
	resp, err := http.Get(fmt.Sprintf(proxyUrl, BEHAVIOR_URL))
	defer resp.Body.Close()

	if err != nil {
		return err, nil
	}

	vconfig := VConfig{}

	rawJSON, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rawJSON, &vconfig)

	return nil, &vconfig
}
