package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	// Fetch leaderboard data
	leaderboardResponse, err := http.Get("https://api.chess.com/pub/leaderboards")
	if err != nil {
		fmt.Println("Error fetching leaderboard data:", err)
		return
	}
	defer leaderboardResponse.Body.Close()
	leaderboardData := make(map[string][]Player)
	err = json.NewDecoder(leaderboardResponse.Body).Decode(&leaderboardData)
	if err != nil {
		fmt.Println("Error decoding leaderboard data:", err)
		return
	}

	// Fetch streamer data
	streamerResponse, err := http.Get("https://api.chess.com/pub/streamers")
	if err != nil {
		fmt.Println("Error fetching streamer data:", err)
		return
	}
	defer streamerResponse.Body.Close()
	var streamerData struct {
		Streamers []Streamer `json:"streamers"`
	}
	err = json.NewDecoder(streamerResponse.Body).Decode(&streamerData)
	if err != nil {
		fmt.Println("Error decoding streamer data:", err)
		return
	}

	// Extract usernames from leaderboard data
	leaderboardUsernames := make(map[string]struct{})
	for _, players := range leaderboardData {
		for _, player := range players {
			leaderboardUsernames[player.Username] = struct{}{}
		}
	}

	// Find and display matches with Twitch URL
	for _, streamer := range streamerData.Streamers {
		if _, exists := leaderboardUsernames[streamer.Username]; exists && streamer.TwitchURL != "" {
			fmt.Printf("Streamer: %s\nTwitch URL: %s\nAvatar: %s\n\n",
				streamer.Username, streamer.TwitchURL, streamer.Avatar)
		}
	}
}
