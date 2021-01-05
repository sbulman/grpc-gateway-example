package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	address := flag.String("server", "http://localhost:8080", "HTTP gateway url")
	flag.Parse()

	t := time.Now().In(time.UTC)
	pfx := t.Format(time.RFC3339Nano)

	var body string

	// Create
	resp, err := http.Post(*address+"/v1/todos", "application/json", strings.NewReader(fmt.Sprintf(`
    {
      "todo": {
        "title":"title (%s)",
        "description":"description (%s)",
        "reminder":"%s"
      }
    }`, pfx, pfx, pfx)))
	if err != nil {
		log.Fatalf("failed to call Create method: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		body = fmt.Sprintf("failed read Create response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	_ = resp.Body.Close()
	log.Printf("Create response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	var created struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(bodyBytes, &created)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON response of Create method: %v", err)
	}

	// Read
	resp, err = http.Get(fmt.Sprintf("%s%s/%s", *address, "/v1/todos", created.ID))
	if err != nil {
		log.Fatalf("failed to call Read method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		body = fmt.Sprintf("failed read Read response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	_ = resp.Body.Close()
	log.Printf("Read response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Update
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", *address, "/v1/todos", created.ID),
		strings.NewReader(fmt.Sprintf(`
		{
			"todo": {
				"title":"title (%s) + updated",
				"description":"description (%s) + updated",
				"reminder":"%s"
			}
		}
	`, pfx, pfx, pfx)))
	if err != nil {
		log.Fatalf("failed to call Update method: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call Update method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		body = fmt.Sprintf("failed read Update response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	_ = resp.Body.Close()
	log.Printf("Update response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// ReadAll
	resp, err = http.Get(*address + "/v1/todos")
	if err != nil {
		log.Fatalf("failed to call ReadAll method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		body = fmt.Sprintf("failed read ReadAll response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	_ = resp.Body.Close()
	log.Printf("ReadAll response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call Delete
	req, err = http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", *address, "/v1/todos", created.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call Delete method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		body = fmt.Sprintf("failed read Delete response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	_ = resp.Body.Close()
	log.Printf("Delete response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}
