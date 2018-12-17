package main

import (
	"encoding/json"
	"net/http"
	"os"
)

const contentTypeHeader, contentType = "Content-Type", "application/json"

func main() {
	http.HandleFunc("/", stan)
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	check(http.ListenAndServe(":"+port, nil))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type entry struct {
	Image string `json:"image"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

func stan(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set(contentTypeHeader, contentType)
	var (
		req struct {
			Payload []struct {
				Drm          bool `json:"drm"`
				EpisodeCount int  `json:"episodeCount"`
				Image        struct {
					ShowImage string `json:"showImage"`
				} `json:"image"`
				Slug  string `json:"slug"`
				Title string `json:"title"`
			} `json:"payload"`
		}
		encoder = json.NewEncoder(w)
	)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		check(encoder.Encode(struct {
			Error string `json:"error"`
		}{"Could not decode request: JSON parsing failed"}))
		return
	}
	var res struct {
		Response []entry `json:"response"`
	}
	for _, x := range req.Payload {
		if x.Drm && x.EpisodeCount > 0 {
			res.Response = append(res.Response, entry{x.Image.ShowImage, x.Slug, x.Title})
		}
	}
	check(encoder.Encode(res))
}
