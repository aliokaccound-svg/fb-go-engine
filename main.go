package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// পোর্ট সেটআপ
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// রিকোয়েস্ট হ্যান্ডেলার
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// পোস্ট রিকোয়েস্ট না হলে হ্যালো দেখাবে
		if r.Method != http.MethodPost {
			fmt.Fprintf(w, "Hello! Railway is running perfectly.")
			return
		}

		// মেসেজ বডি পড়া
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		// লগ প্রিন্ট করা (Railway Logs-এ দেখার জন্য)
		log.Printf("📩 New Message: %s", string(body))
		
		// সাকসেস রেসপন্স দেওয়া
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Received"))
	})

	// সার্ভার স্টার্ট
	log.Println("Server starting on 0.0.0.0:" + port)
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal(err)
	}
}
