package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func makeDefaultWebhook(
	repoName string,
	codename string,
	arch string,
) func(string) {
	type Repository struct {
		Version  string `json:"version"`
		Codename string `json:"codename"`
		Arch     string `json:"arch"`
	}
	type ClientPayload struct {
		Repository Repository `json:"repository"`
	}
	type WebhookPayloadType struct {
		EventType     string        `json:"event_type"`
		ClientPayload ClientPayload `json:"client_payload"`
	}

	return func(version string) {
		log.Println("webhook for " + version + " to " + repoName)

		webhookPayload := WebhookPayloadType{
			EventType: "trigger_build",
			ClientPayload: ClientPayload{
				Repository: Repository{
					Codename: codename,
					Version:  version,
					Arch:     arch,
				},
			},
		}

		webhookPayloadJson, err := json.Marshal(webhookPayload)
		if err != nil {
			panic(err)
		}

		log.Println(string(webhookPayloadJson))

		httpReq, err := http.NewRequest(http.MethodPost, "https://api.github.com/repos/"+repoName+"/dispatches", bytes.NewBuffer(webhookPayloadJson))
		if err != nil {
			panic(err)
		}

		token := os.Getenv("GH_TOKEN_TRIGGER")
		if token == "" {
			panic(fmt.Errorf("missing token in GH_TOKEN_TRIGGER"))
		}
		httpReq.Header.Set("Authorization", "Bearer "+token)
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		io.Copy(ioutil.Discard, resp.Body)

		log.Println("webhook result:", resp.StatusCode)
	}
}
