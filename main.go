package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func fetchDataFromResponse(req *http.Request) {
	// Perform Http request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Image data:", len(imageData))
}

func RequestWithoutContext() {
	// Create Http request
	req, err := http.NewRequest("GET", "https://golang.org/doc/gopher/frontpage.png", nil)
	if err != nil {
		panic(err)
	}

	// Perform Http request
	fetchDataFromResponse(req)
}

func RequestWithTimeoutContext(ctx context.Context) {
	// Create Http request
	req, err := http.NewRequestWithContext(ctx, "GET", "https://golang.org/doc/gopher/frontpage.png", nil)
	if err != nil {
		panic(err)
	}

	// Perform Http request
	fetchDataFromResponse(req)
}

func RequestUsingGin() {
	// Example 3 - Request using Gin framework with standard context
	r := gin.Default()
	r.GET("/hello", func(ctx *gin.Context) {
		// Here with ctx.Request.Context() we are creating a context from request context. In case of request cancellation, it will cancel the request.
		// Error - Get "https://golang.org/doc/gopher/frontpage.png": context canceled
		req, err := http.NewRequestWithContext(ctx.Request.Context(), "GET", "https://golang.org/doc/gopher/frontpage.png", nil)
		if err != nil {
			panic(err)
		}

		// Perform Http request
		fetchDataFromResponse(req)
	}, nil)

	// Example 4 - Request using Gin framework with timeout context
	r.GET("/hello", func(ctx *gin.Context) {
		timeoutContext, cancel := context.WithTimeout(ctx.Request.Context(), 100*time.Millisecond)
		defer cancel()

		// Here if request is not completed in 100 millisecond then it will cancel the request and return error.
		// Error = Get "https://golang.org/doc/gopher/frontpage.png": context deadline exceeded
		req, err := http.NewRequestWithContext(timeoutContext, "GET", "https://golang.org/doc/gopher/frontpage.png", nil)
		if err != nil {
			panic(err)
		}
		// Perform Http request
		fetchDataFromResponse(req)
	}, nil)

}

func main() {
	// Context - It is used to pass request scoped data, cancelation signal, deadlines to all goroutines involved in handling a request.
	// Context is like a tree. You can create a context from another context and it will inherit all the values from parent context.

	// Example 1 - Request without context
	// RequestWithoutContext()

	// Example 2 - Request with context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond) // Context.Background() is used to create a root context
	defer cancel()
	// Here if request is not completed in 100 millisecond then it will cancel the request and return error.
	// Error = Get "https://golang.org/doc/gopher/frontpage.png": context deadline exceeded
	RequestWithTimeoutContext(timeoutContext)

	// Example 3 - Request using Gin framework
	RequestUsingGin()

}
