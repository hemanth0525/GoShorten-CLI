# Go URL Shortener

## What’s This?

This is a simple URL shortener project built with Go. It’s a great beginner project to get hands-on with Go. You’ll learn how to create a basic web server, use goroutines to run tasks in the background, and handle data with maps.

## What You’ll Learn

- **Go Modules**: How to set up and manage your Go project.
- **HTTP Servers**: How to build a web server with Go.
- **Goroutines**: Running the server and CLI in parallel.
- **Maps and Concurrency**: Storing and accessing data safely.

## Project Structure

You only need one file for this project:

```
GoShorten-CLI/
│
├── main.go
├── go.mod
└── README.md
```

### **`main.go`**

This file contains all the code to handle URL shortening and redirection.

#### Code Breakdown

```go
package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// URLStore is a map that links short URLs to long URLs.
var (
	URLStore = make(map[string]string)
	mu       sync.Mutex
)

// RedirectHandler redirects short URLs to the long URL.
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:] // Get short URL from the request
	mu.Lock()
	longURL, exists := URLStore[shortURL]
	mu.Unlock()

	if exists {
		http.Redirect(w, r, longURL, http.StatusFound)
	} else {
		http.NotFound(w, r) // Show 404 if the short URL isn’t found
	}
}

func main() {
	// Start the HTTP server in a separate goroutine.
	go func() {
		http.HandleFunc("/", RedirectHandler)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Server failed to start: %v\n", err)
		}
	}()

	fmt.Println("Server running at http://localhost:8080")

	for {
		var longURL, shortURL string

		// Get the long URL from the user.
		fmt.Print("Enter long URL: ")
		fmt.Scanln(&longURL)

		// Add "https://" if the URL doesn’t start with it.
		if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
			longURL = "https://" + longURL
		}

		// Get a custom short URL or generate one if none is provided.
		fmt.Print("Enter custom short URL (optional): ")
		fmt.Scanln(&shortURL)
		if shortURL == "" {
			shortURL = fmt.Sprintf("%d", len(URLStore)+1)
		}

		// Save the mapping from short URL to long URL.
		mu.Lock()
		URLStore[shortURL] = longURL
		mu.Unlock()

		fmt.Printf("Short URL: http://localhost:8080/%s\n\n", shortURL)
	}
}
```

#### What’s Happening?

1. **Imports**: Includes libraries for HTTP, strings, and concurrency.
2. **Global Variables**:
   - `URLStore`: Keeps track of short URLs and their long counterparts.
   - `mu`: A mutex to ensure thread safety when multiple users are interacting with the URL shortener.
3. **RedirectHandler**: Handles redirection from short URLs to long URLs.
4. **Main Function**:
   - Starts the HTTP server in the background.
   - Continuously prompts for URLs and handles user input.
   - Stores the URLs and prints out the shortened URL.

## How to Use

### What You Need

- Go 1.16 or later installed on your computer.

### Steps to Run It

1. **Clone the Repo**:
   ```bash
   git clone https://github.com/hemanth0525/GoShorten-CLI.git
   ```

2. **Set Up Go Modules**:
   ```bash
   go mod init GoShorten-CLI
   go mod tidy
   ```

3. **Run the Code**:
   ```bash
   go run main.go
   ```
   or (if you have `go.mod` in your project):
   ```bash
   go run .
   ```

### How to Use It

1. Enter a long URL when prompted (e.g., `google.com`).
2. Provide a custom short URL or just press Enter to let the program create one for you.
3. The program will give you a short URL that you can use to redirect to the original URL.

### Example

```plaintext
Enter long URL:
google.com
Enter custom short URL (optional):
goog
Short URL: http://localhost:8080/goog
```

Visit `http://localhost:8080/goog` in your browser to go to `https://google.com`.

## License

This project is licensed under the MIT License. Check out the [LICENSE](LICENSE) file for details.