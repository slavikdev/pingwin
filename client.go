package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const DROP_ERROR_RATE = 0.5

type Client struct {
	id        int
	targetUrl string
}

func NewClient(id int, targetHost string) *Client {
	return &Client{id: id, targetUrl: targetHost}
}

func (client *Client) Run() {
	numRequests := 0
	numFailures := 0
	httpClient := &http.Client{}

	for {
		req := client.CreateRequest()
		if req == nil {
			numFailures++
		} else {
			resp, err := httpClient.Do(client.CreateRequest())
			if err != nil {
				numFailures++
			} else {
				defer resp.Body.Close()
				_, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					numFailures++
				}
			}
		}

		numRequests++

		errorRate := float64(numFailures) / float64(numRequests)
		if numRequests%1000 == 0 {
			client.Log("%d requests completed, %d failures (error rate = %f)", numRequests, numFailures, errorRate)
			if errorRate > DROP_ERROR_RATE {
				client.Log("Error rate too high, aborting the task.")
				break
			}
		}
	}
}

func (client *Client) CreateRequest() *http.Request {
	req, err := http.NewRequest("GET", client.targetUrl, nil)
	if err != nil {
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "Keep-Alive")

	return req
}

func (client *Client) Log(message string, args ...interface{}) {
	fmt.Printf("[CL-%d, %s] %s\n", client.id, client.targetUrl, fmt.Sprintf(message, args...))
}
