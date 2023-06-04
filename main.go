// https://goalice-1-f0793019.deta.app/
package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"net/http"
	"os"
)

type ResponseToUser struct {
	Response Response `json:"response"`
	Version  string   `json:"version"`
}

type Response struct {
	Text       string `json:"text"`
	Tts        string `json:"tts"`
	EndSession bool   `json:"end_session"`
}

type UserRequest struct {
	Session string  `json:"session"`
	Request Request `json:"request"`
	Version string  `json:"version"`
}

type Request struct {
	Type              string `json:"type"`
	Command           string `json:"command"`
	OriginalUtterance string `json:"original_utterance"`
	Nlu               Nlu    `json:"nlu"`
}

type Nlu struct {
	Tokens []string `json:"tokens"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handler)

	log.Printf("App listening on port %s!", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	} else {
		log.Printf(string(body))
	}
	// decoder := json.NewDecoder(r.Body)
	req := UserRequest{}
	json.Unmarshal([]byte(body), &req)

	response_text := "Я пока не знаю ответа на этот вопрос"
	if req.Request.OriginalUtterance == "" {
		response_text = "Привет!"
	} else if slices.Contains(req.Request.Nlu.Tokens, "биткоина") {
		response_text = "Как обычно очень высокий :)"
	}

	var response ResponseToUser
	response.Response.Text = response_text
	response.Response.EndSession = false
	response.Version = req.Version
	response_json, err := json.Marshal(response)
	// e:= decoder.Decode(&obj)
	log.Printf(string(response_json))

	fmt.Fprintf(w, string(response_json))

}
