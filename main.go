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
	
	Attrs  map[string]string `json:"-"`
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
		// ScreenSize:      req.Sc,
		// Attributes: map[string]attribute{
		// 	req.Atrk1: {
		// 		Value: req.Atrv1,
		// 		Type:  req.Atrt1,
		// 	},
		// 	req.Atrk2: {
		// 		Value: req.Atrv2,
		// 		Type:  req.Atrt2,
		// 	},
		// },
		// Traits: map[string]trait{
		// 	req.Uatrk1: {
		// 		Value: req.Uatrv1,
		// 		Type:  req.Uatrt1,
		// 	},
		// 	req.Uatrk2: {
		// 		Value: req.Uatrv2,
		// 		Type:  req.Uatrt2,
		// 	},
		// 	req.Uatrv3: {
		// 		Value: req.Uatrv3,
		// 		Type:  req.Uatrt3,
		// 	},
		// },
	}
	for k, v := range req.Attrs {
		if k[:4] == "atrk" {
			resp.Attributes[v] = attribute{
				Value: req.Attrs[k[0:3] + "v" + k[4:]],
				Type:  req.Attrs[k[0:3] + "t" + k[4:]],
			}
		}else if k[:5] == "uatrk" {
			resp.Traits[v] = trait{
				Value: req.Attrs[k[0:4] + "v" + k[5:]],
				Type:  req.Attrs[k[0:4] + "t" + k[5:]],
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	wg.Done()
}
