package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Notification request headers
const (
	TWITCH_MESSAGE_ID        = "twitch-eventsub-message-id"
	TWITCH_MESSAGE_TIMESTAMP = "twitch-eventsub-message-timestamp"
	TWITCH_MESSAGE_SIGNATURE = "twitch-eventsub-message-signature"
	MESSAGE_TYPE             = "twitch-eventsub-message-type"

	// Notification message types
	MESSAGE_TYPE_VERIFICATION = "webhook_callback_verification"
	MESSAGE_TYPE_NOTIFICATION = "notification"
	MESSAGE_TYPE_REVOCATION   = "revocation"

	// Prepend this string to the HMAC that's created from the message
	HMAC_PREFIX = "sha256="
)

func main() {
	http.HandleFunc("/eventsub", eventSubHandler)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func eventSubHandler(w http.ResponseWriter, r *http.Request) {
	// Read the body once and use it throughout
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error - Unable to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Prepare HMAC for verification
	secret := getSecret()
	message := r.Header.Get(TWITCH_MESSAGE_ID) + r.Header.Get(TWITCH_MESSAGE_TIMESTAMP) + string(body)
	hmac := HMAC_PREFIX + getHmac(secret, message)

	// Verify HMAC Signature
	if !verifyMessage(hmac, r.Header.Get(TWITCH_MESSAGE_SIGNATURE)) {
		fmt.Println("403 - Signatures didn't match.")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Handle the verified message
	messageType := r.Header.Get(MESSAGE_TYPE)
	var notification map[string]interface{}
	if err := json.Unmarshal(body, &notification); err != nil {
		http.Error(w, "Internal Server Error - Unable to parse JSON", http.StatusInternalServerError)
		return
	}

	switch messageType {
	case MESSAGE_TYPE_NOTIFICATION:
		fmt.Printf("Event type: %s\n", notification["subscription"].(map[string]interface{})["type"])
		fmt.Printf("%s\n", string(body))
		w.WriteHeader(http.StatusNoContent)
	case MESSAGE_TYPE_VERIFICATION:
		challenge := notification["challenge"].(string)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
	case MESSAGE_TYPE_REVOCATION:
		w.WriteHeader(http.StatusNoContent)
		fmt.Printf("%s notifications revoked!\n", notification["subscription"].(map[string]interface{})["type"])
		fmt.Printf("reason: %s\n", notification["subscription"].(map[string]interface{})["status"])
		fmt.Printf("condition: %s\n", notification["subscription"].(map[string]interface{})["condition"])
	default:
		w.WriteHeader(http.StatusNoContent)
		fmt.Printf("Unknown message type: %s\n", messageType)
	}
}

func getSecret() string {
	// TODO: Get secret from secure storage. This is the secret you pass
	// when you subscribed to the event.
	return "520q79qc0ai4c4kpe1yzt7ipxl5jkc"
}

// Get the HMAC.
func getHmac(secret, message string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// Verify whether our hash matches the hash that Twitch passed in the header.
func verifyMessage(hmac, verifySignature string) bool {
	return hmac == verifySignature
}
