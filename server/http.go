package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type httpServer struct{}

// Check to see if an item is recycable. Checks the param "item".
func (s *httpServer) recyclable(w http.ResponseWriter, r *http.Request) {
	material := r.URL.Query().Get("item")
	fmt.Fprintf(w, "We received your query regarding %s", material)
}

// Used to received requests from our frontend containing user location info.
type RequestInfo struct {
	Material  string  `json:"material"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Used to send requests to external locations service.
type LocServiceRequest struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

// The response that we receive from our locations service.
type LocationReponse struct {
	ErrorMsg    string `json:"errorMessage"`
	City        string `json:"city"`
	State       string `json:"state"`
	CommunityID string `json:"communityID"`
}

// Used to send request to guideliens service
type GuidelinesRequest struct {
	CommunityID string `json:"communityID"`
}

// Result from individual material in GuidelinesResponse
type MaterialResult struct {
	MID          int32  `json:"mID"`
	CommunityID  string `json:"communityID"`
	Category     string `json:"category"`
	YesNo        string `json:"yesNo"`
	CategoryType string `json:"categoryType"`
	Material     string `json:"material"`
}

// Values sent back from request to recycling service
type GuidelinesResponse struct {
	Guidelines []MaterialResult `json:"guidelines"`
}

func (s *httpServer) location(w http.ResponseWriter, r *http.Request) {
	// Provides lat and long for first service and material for second service
	var reqInfo RequestInfo
	err := json.NewDecoder(r.Body).Decode(&reqInfo)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Send request to location service
	body, err := json.Marshal(
		&LocServiceRequest{
			reqInfo.Latitude,
			reqInfo.Longitude},
	)
	res, err := http.Post("http://localhost:8081", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var locRes LocationReponse
	err = json.NewDecoder(res.Body).Decode(&locRes)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if locRes.ErrorMsg != "" {
		log.Println(locRes.ErrorMsg)
		http.Error(w, locRes.ErrorMsg, http.StatusNotFound)
	}
	// Send request to recycling info service
	body, err = json.Marshal(
		&GuidelinesRequest{CommunityID: locRes.CommunityID},
	)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err = http.Post("http://localhost:8082", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var guideResults GuidelinesResponse
	err = json.NewDecoder(res.Body).Decode(&guideResults)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if locRes.ErrorMsg != "" {
		log.Println(locRes.ErrorMsg)
		http.Error(w, locRes.ErrorMsg, http.StatusNotFound)
	}
	// Send back results
	for _, guideline := range guideResults.Guidelines {
		if guideline.Category == reqInfo.Material {
			fmt.Fprintf(w, "The answer for %s is %s", reqInfo.Material, guideline.YesNo)
			break
		}
	}
}

// Returns an *http.Server that listens at location `addr` with routes registered.
func NewHTTPServer(addr string) *http.Server {
	httpServer := &httpServer{}

	r := mux.NewRouter()
	r.HandleFunc("/location", httpServer.location).Methods("POST")
	r.HandleFunc("/recycle", httpServer.recyclable).Methods("POST")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
