package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/michiwend/gomusicbrainz"
)

func main() {
	// customize output
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	//read args
	if len(os.Args) < 4 {
		log.Fatalf(`Usage: ./crud <description> <version> <contact>`)
	}
	description, version, contact := os.Args[1], os.Args[2], os.Args[3]

	//create client
	client, err := gomusicbrainz.NewWS2Client("https://musicbrainz.org/ws/2", description, version, contact)
	if err != nil {
		log.Fatalf("Failed to create WS2 client: %v", err)
	}

	handler, err := newArtistHandler(client)
	if err != nil {
		log.Fatalf("Failed to create /artist handler: %v", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/artist/search/{name}", handler.searchArtist)
	r.HandleFunc("/artist/lookup/{id}", handler.lookupArtist)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start http server: %v", err)
	}
}

type artistHandler struct {
	client   *gomusicbrainz.WS2Client
	template *template.Template
}

func newArtistHandler(client *gomusicbrainz.WS2Client) (*artistHandler, error) {
	funcMap := template.FuncMap{"formatTime": formatTime}
	tmpl, err := template.New("base").Funcs(funcMap).ParseFiles("artists.html", "artist.html")
	if err != nil {
		return nil, fmt.Errorf("Failed to parse artists.html template: %v", err)
	}
	return &artistHandler{client: client, template: tmpl}, nil
}

func (h *artistHandler) searchArtist(w http.ResponseWriter, r *http.Request) {
	artistName := mux.Vars(r)["name"]
	if artistName == "" {
		writeError(w, http.StatusBadRequest, "Expecting non-empty artist name")
		return
	}

	artistResponse, err := h.client.SearchArtist(artistName, -1, -1)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Failed to find artist: %v", err))
		return
	}

	if err := h.template.ExecuteTemplate(w, "artists.html", artistResponse); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to execute artists template: %v", err))
	}
	return
}

func (h *artistHandler) lookupArtist(w http.ResponseWriter, r *http.Request) {
	artistID := mux.Vars(r)["id"]
	if artistID == "" {
		writeError(w, http.StatusBadRequest, "Expecting non-empty artist ID")
		return
	}

	artistResponse, err := h.client.LookupArtist(gomusicbrainz.MBID(artistID))
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Failed to lookup artist: %v", err))
		return
	}

	if err := h.template.ExecuteTemplate(w, "artist.html", artistResponse); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to execute artist template: %v", err))
	}
	return
}

func writeError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}

func formatTime(input time.Time) string {
	if input.IsZero() {
		return ""
	}
	return input.Format("2006-01-02")
}
