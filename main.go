package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// URLStore keeps track of short to long URL mappings.
var (
	URLStore = make(map[string]string)
	mu       sync.Mutex
)

// RedirectHandler is called when someone visits a short URL.
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:] // Get the part after the "/"
	mu.Lock()
	longURL, ok := URLStore[shortURL]
	mu.Unlock()

	if ok {
		http.Redirect(w, r, longURL, http.StatusFound)
	} else {
		http.NotFound(w, r) // Oops! Short URL not found.
	}
}

func main() {
	// Start the HTTP server in the background.
	go func() {
		http.HandleFunc("/", RedirectHandler)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Oops! Server didn't start: %v\n", err)
		}
	}()

	fmt.Println("Server is up and running at http://localhost:8080")

	for {
		var longURL, shortURL string

		// Ask for the long URL.
		fmt.Print("Enter long URL: ")
		fmt.Scanln(&longURL)

		// Add "https://" if the URL doesn’t start with "http://" or "https://"
		if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
			longURL = "https://" + longURL
		}

		// Ask for a custom short URL.
		fmt.Print("Enter custom short URL (optional): ")
		fmt.Scanln(&shortURL)

		// If no custom URL is given, generate one.
		if shortURL == "" {
			shortURL = fmt.Sprintf("%d", len(URLStore)+1)
		}

		// Store the URL mapping.
		mu.Lock()
		URLStore[shortURL] = longURL
		mu.Unlock()

		fmt.Printf("Your short URL is: http://localhost:8080/%s\n\n", shortURL)
	}
}