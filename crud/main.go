package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/michiwend/gomusicbrainz"
)

var client *gomusicbrainz.WS2Client

func main() {
	var err error
	client, err = gomusicbrainz.NewWS2Client("https://musicbrainz.org/ws/2", "My personal MB feed service", "0.0.1", "https://github.com/n0vice")
	if err != nil {
		log.Fatalf("Failed to create WS2 client: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/artist/{name}", artistHandler)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start http server: %v", err)
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artistName := vars["name"]

	artistResponse, err := client.SearchArtist(artistName, -1, -1)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Failed to find artist: %v", err)
		return
	}

	for _, a := range artistResponse.Artists {
		fmt.Fprintln(w, artist(*a))
	}
	return
}

type artist gomusicbrainz.Artist

func (a artist) String() string {
	result := a.Name
	var additionalInfo string
	var t time.Time
	if a.Lifespan.Begin.Time != t {
		additionalInfo += fmt.Sprintf("%v", a.Lifespan.Begin.Year())
		if a.Lifespan.Ended {
			additionalInfo += fmt.Sprintf("-%v", a.Lifespan.End.Year())
		}
		additionalInfo += " "
	}
	if a.Area.Name != "" {
		additionalInfo += a.Area.Name
		if a.BeginArea.Name != "" {
			additionalInfo += ", " + a.BeginArea.Name
		}
	}
	if additionalInfo != "" {
		result += fmt.Sprintf(" (%v)", additionalInfo)
	}
	if len(a.Tags) > 0 {
		var index, count int
		for i, t := range a.Tags {
			if t.Count > count {
				count = t.Count
				index = i
			}
		}
		result += a.Tags[index].Name
	}
	return result
}
