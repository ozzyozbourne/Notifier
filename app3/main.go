package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Define the structure for the WhatsApp message payload
type WhatsAppMessage struct {
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             Text   `json:"text"`
}

type Text struct {
	Body string `json:"body"`
}

func sendWhatsAppMessage(w http.ResponseWriter, r *http.Request) {
	accessToken := "YOUR_ACCESS_TOKEN"
	phoneNumberId := "YOUR_PHONE_NUMBER_ID"
	recipientPhoneNumber := "+1234567890" // Replace with the recipient's phone number

	// Parse request body (assuming JSON input)
	var inputMessage struct {
		Message string `json:"message"`
	}

	err := json.NewDecoder(r.Body).Decode(&inputMessage)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create the WhatsAppMessage payload
	message := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		To:               recipientPhoneNumber,
		Type:             "text",
		Text: Text{
			Body: inputMessage.Message,
		},
	}

	messageData, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Failed to serialize message", http.StatusInternalServerError)
		return
	}

	// Send the request to the WhatsApp Business API
	req, err := http.NewRequest("POST", fmt.Sprintf("https://graph.facebook.com/v13.0/%s/messages", phoneNumberId), bytes.NewBuffer(messageData))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body) // Use io.ReadAll instead of ioutil.ReadAll
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, string(respBody), resp.StatusCode)
		return
	}

	fmt.Fprintln(w, "Message sent successfully")
}

func main() {
	http.HandleFunc("/send-whatsapp-message", sendWhatsAppMessage)

	fmt.Println("Starting server on :8080 with TLS...")
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
