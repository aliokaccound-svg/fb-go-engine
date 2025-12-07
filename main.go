package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// ১. ফেসবুকের ডাটা স্ট্রাকচার (JSON Model)
// ফেসবুক যে প্যাকেট পাঠায়, সেটা রিসিভ করার জন্য এই ছাঁচ (Struct) লাগবে
type WebhookEvent struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Message struct {
				Text string `json:"text"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

func main() {
	// পোর্ট সেটআপ
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// শুধুমাত্র POST মেথড এলাউড
		if r.Method != http.MethodPost {
			fmt.Fprintf(w, "Bot Engine is Running! Send a POST request.")
			return
		}

		// ২. বডি রিড করা
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		// ৩. জেসন পার্স করা (JSON to Go Struct)
		var event WebhookEvent
		err = json.Unmarshal(body, &event)
		if err != nil {
			log.Printf("Could not parse JSON: %v", err)
			// ফেইল করলেও আমরা 200 দেব, নাহলে ফেসবুক ব্লক করবে
			w.WriteHeader(http.StatusOK) 
			return
		}

		// ৪. লুপ চালিয়ে মেসেজ বের করা
		// ফেসবুক মাঝে মাঝে একসাথে অনেকগুলো মেসেজ পাঠায় (Batch), তাই লুপ লাগবে
		for _, entry := range event.Entry {
			for _, messaging := range entry.Messaging {
				// শুধু টেক্সট মেসেজ আসলে কাজ করব (ছবি/স্টিকার আপাতত বাদ)
				if messaging.Message.Text != "" {
					senderID := messaging.Sender.ID
					userMessage := messaging.Message.Text

					// লগে প্রিন্ট করছি (এখানেই পরে AI বা ডাটাবেস বসবে)
					log.Printf("✅ User [%s] said: %s", senderID, userMessage)
					
					// TODO: এখানেই Supabase এ সেভ করার কোড বসবে
				}
			}
		}

		// ৫. সাকসেস রেসপন্স
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Event Processed"))
	})

	log.Println("Server starting on 0.0.0.0:" + port)
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal(err)
	}
}
