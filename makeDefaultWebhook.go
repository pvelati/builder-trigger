package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func makeDefaultWebhook(
	repoName string,
	codename string,
	arch string,
	build_arch string,
) func(string) {
	type Repository struct {
		Version    string `json:"version"`
		Codename   string `json:"codename"`
		Arch       string `json:"arch"`
		Build_arch string `json:"build_arch"`
	}
	type ClientPayload struct {
		Repository Repository `json:"repository"`
	}
	type WebhookPayloadType struct {
		EventType     string        `json:"event_type"`
		ClientPayload ClientPayload `json:"client_payload"`
	}

	return func(version string) {
		log.Println("skip webhook")
		return
		log.Println("webhook for " + version + " to " + repoName)

		webhookPayload := WebhookPayloadType{
			EventType: "trigger_build",
			ClientPayload: ClientPayload{
				Repository: Repository{
					Codename:   codename,
					Version:    version,
					Arch:       arch,
					Build_arch: build_arch,
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

		httpReq.Header.Set("Authorization", "Bearer "+githubToken())
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		io.Copy(io.Discard, resp.Body)

		if resp.StatusCode != 200 {
			panic(fmt.Errorf("invalid status %s", resp.Status))
		}

		log.Println("webhook result:", resp.StatusCode)
	}
}
