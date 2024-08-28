package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Player struct {
	Username string `json:"username"`
}

type Streamer struct {
	Username  string `json:"username"`
	TwitchURL string `json:"twitch_url"`
	Avatar    string `json:"avatar"`
}

type PageData struct {
	Streamers []Streamer
}

func main() {
	http.HandleFunc("/", displayStreamers)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func displayStreamers(w http.ResponseWriter, r *http.Request) {
	// Fetch leaderboard data
	leaderboardResponse, err := http.Get("https://api.chess.com/pub/leaderboards")
	if err != nil {
		http.Error(w, "Failed to fetch leaderboard data", http.StatusInternalServerError)
		return
	}
	defer leaderboardResponse.Body.Close()
	leaderboardData := make(map[string][]Player)
	err = json.NewDecoder(leaderboardResponse.Body).Decode(&leaderboardData)
	if err != nil {
		http.Error(w, "Failed to decode leaderboard data", http.StatusInternalServerError)
		return
	}

	// Fetch streamer data
	streamerResponse, err := http.Get("https://api.chess.com/pub/streamers")
	if err != nil {
		http.Error(w, "Failed to fetch streamer data", http.StatusInternalServerError)
		return
	}
	defer streamerResponse.Body.Close()
	var streamerData struct {
		Streamers []Streamer `json:"streamers"`
	}
	err = json.NewDecoder(streamerResponse.Body).Decode(&streamerData)
	if err != nil {
		http.Error(w, "Failed to decode streamer data", http.StatusInternalServerError)
		return
	}

	// Extract usernames from leaderboard data
	leaderboardUsernames := make(map[string]struct{})
	for _, players := range leaderboardData {
		for _, player := range players {
			leaderboardUsernames[player.Username] = struct{}{}
		}
	}

	// Filter streamers with Twitch URL and match them with leaderboard usernames
	var matchedStreamers []Streamer
	for _, streamer := range streamerData.Streamers {
		if _, exists := leaderboardUsernames[streamer.Username]; exists && streamer.TwitchURL != "" {
			matchedStreamers = append(matchedStreamers, streamer)
		}
	}

	// Serve the HTML template
	tmpl := template.Must(template.ParseFiles("template.html"))
	err = tmpl.Execute(w, PageData{Streamers: matchedStreamers})
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
