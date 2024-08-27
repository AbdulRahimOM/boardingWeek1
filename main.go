package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
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

	Atrk1 string `json:"atrk1"`
	Atrv1 string `json:"atrv1"`
	Atrt1 string `json:"atrt1"`

	Atrk2 string `json:"atrk2"`
	Atrv2 string `json:"atrv2"`
	Atrt2 string `json:"atrt2"`

	Uatrk1 string `json:"uatrk1"`
	Uatrv1 string `json:"uatrv1"`
	Uatrt1 string `json:"uatrt1"`

	Uatrk2 string `json:"uatrk2"`
	Uatrv2 string `json:"uatrv2"`
	Uatrt2 string `json:"uatrt2"`

	Uatrk3 string `json:"uatrk3"`
	Uatrv3 string `json:"uatrv3"`
	Uatrt3 string `json:"uatrt3"`
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

type trait struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}
type attribute struct {
	Value string `json:"value"`
	Type  string `json:"type"`
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

	wg := sync.WaitGroup{}
	wg.Add(1)
	go worker(req, &wg, w)
	wg.Wait()
}

func main() {
	http.HandleFunc("/process", helloPostHandler)

	fmt.Printf("Starting server at port 8099\n")
	if err := http.ListenAndServe(":8099", nil); err != nil {
		log.Fatal(err)
	}
}

func worker(req request, wg *sync.WaitGroup, w http.ResponseWriter) {
	resp:= response{
		Event:           req.Ev,
		EventType:       req.Et,
		AppId:           req.Id,
		UserId:          req.Uid,
		MessageId:       req.Mid,
		PageTitle:       req.T,
		PageUrl:         req.P,
		BrowserLanguage: req.L,
		ScreenSize:      req.Sc,
		Attributes: map[string]attribute{
			req.Atrk1: {
				Value: req.Atrv1,
				Type:  req.Atrt1,
			},
			req.Atrk2: {
				Value: req.Atrv2,
				Type:  req.Atrt2,
			},
		},
		Traits: map[string]trait{
			req.Uatrk1: {
				Value: req.Uatrv1,
				Type:  req.Uatrt1,
			},
			req.Uatrk2: {
				Value: req.Uatrv2,
				Type:  req.Uatrt2,
			},
			req.Uatrv3: {
				Value: req.Uatrv3,
				Type:  req.Uatrt3,
			},
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	wg.Done()
}
