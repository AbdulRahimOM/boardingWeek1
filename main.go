package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type request struct {
	Ev  string `json:"ev"`
	Et  string `json:"et"`
	Id  string `json:"id"`
	Uid string `json:"uid"`
	Mid string `json:"mid"`
	T   string `json:"t"`
	P   string `json:"p"`
	L   string `json:"l"`
	Sc  string `json:"sc"`

	CustomFields map[string]string `json:"-"`
}

type response struct {
	Event           string               `json:"event"`
	EventType       string               `json:"event_type"`
	AppId           string               `json:"app_id"`
	UserId          string               `json:"user_id"`
	MessageId       string               `json:"message_id"`
	PageTitle       string               `json:"page_title"`
	PageUrl         string               `json:"page_url"`
	BrowserLanguage string               `json:"browser_language"`
	ScreenSize      string               `json:"screen_size"`
	Attributes      map[string]attribute `json:"attributes"`
	Traits          map[string]trait     `json:"traits"`
}

type attribute struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}
type trait struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

var requestChannel = make(chan request)
var responseChannel = make(chan response)

func main() {
	go worker()
	http.HandleFunc("/process", helloPostHandler)
	fmt.Printf("Starting server at port 8099\n")
	if err := http.ListenAndServe(":8099", nil); err != nil {
		log.Fatal(err)
	}
}

func (r *request) UnmarshalJSON(data []byte) error {
	type Alias request
	if err := json.Unmarshal(data, (*Alias)(r)); err != nil {
		return err
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}
	r.CustomFields = make(map[string]string)
	for key, value := range rawMap {
		switch key {
		case "ev", "et", "id", "uid", "mid", "t", "p", "l", "sc":
			continue
		}

		var fieldValue string
		if err := json.Unmarshal(value, &fieldValue); err != nil {
			fieldValue = string(value) // Fallback to raw string if not a simple JSON string
		}
		r.CustomFields[key] = fieldValue
	}

	return nil
}

func helloPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	requestChannel <- req

	resp := <-responseChannel
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)

}

// func worker(req request, wg *sync.WaitGroup, w http.ResponseWriter) {
func worker() {
	for {
		req := <-requestChannel
		resp := response{
			Event:           req.Ev,
			EventType:       req.Et,
			AppId:           req.Id,
			UserId:          req.Uid,
			MessageId:       req.Mid,
			PageTitle:       req.T,
			PageUrl:         req.P,
			BrowserLanguage: req.L,
			ScreenSize:      req.Sc,
			Attributes:      make(map[string]attribute),
			Traits:          make(map[string]trait),
		}

		for k, v := range req.CustomFields {
			if k[:4] == "atrk" {
				resp.Attributes[v] = attribute{
					Value: req.CustomFields[k[0:3]+"v"+k[4:]],
					Type:  req.CustomFields[k[0:3]+"t"+k[4:]],
				}
			} else if k[:5] == "uatrk" {
				resp.Traits[v] = trait{
					Value: req.CustomFields[k[0:4]+"v"+k[5:]],
					Type:  req.CustomFields[k[0:4]+"t"+k[5:]],
				}
			}
		}
		responseChannel <- resp
	}

}
