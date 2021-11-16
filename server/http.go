package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gorilla/mux"
)

type httpServer struct{}

func (s *httpServer) home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.FromSlash("./ui/html/home.page.tmpl"),
		filepath.FromSlash("./ui/html/base.layout.tmpl"),
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

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

func (s *httpServer) location(w http.ResponseWriter, r *http.Request) {
	log.Println("HIT")
	var reqInfo RequestInfo
	err := json.NewDecoder(r.Body).Decode(&reqInfo)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// fmt.Fprintf(w, "<h1>Your lat and long are %2f and %2f respectively</h1>", locInfo.Latitude, locInfo.Longitude)
	body, err := json.Marshal(
		&LocServiceRequest{
			reqInfo.Latitude,
			reqInfo.Longitude},
	)
	res, err := http.Post("http://localhost:8080", "application/json", bytes.NewBuffer(body))
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
	fmt.Fprintf(w, "You live in %s, %s", locRes.City, locRes.State)
}

// Returns an *http.Server that listens at location `addr` with routes registered.
func NewHTTPServer(addr string) *http.Server {
	httpServer := &httpServer{}

	r := mux.NewRouter()
	r.HandleFunc("/", httpServer.home).Methods("GET")
	// r.HandleFunc("/", httpServer.recyclable).Methods("GET")
	r.HandleFunc("/location", httpServer.location).Methods("POST")
	r.HandleFunc("/recycle", httpServer.recyclable).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
