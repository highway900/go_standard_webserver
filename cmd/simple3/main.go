package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

const HeaderRequestID string = "X-Request-Id"

// ContextKey is a key for the context to lookup data
type ContextKey string

const contextKeyRequestID ContextKey = "requestID"

// Request ID is an atomic counter
var requestIDCounter uint64

// User struct defines a user
type User struct {
	ID        int     `json:"ID"`
	Balance   float32 `json:"Balance"`
	accountID string  // automatically omitted from JSON because it is hidden/private
}

// really you should use a UUID generated on demand
func getRequestID(ctx context.Context) uint64 {
	reqIDRaw := ctx.Value(contextKeyRequestID) // reqIDRaw at this point is of type 'interface{}'
	reqID, ok := reqIDRaw.(uint64)             // be careful to cast correctly here
	if !ok {
		return 0
	}
	return reqID
}

// LogError is a log wrapper function for printing errors
func LogError(msg string, err error) {
	log.Println(fmt.Errorf("ERROR - %s: %v", msg, err))
}

/////////////////////////////////////////////////////////////////////
// Middleware
/////////////////////////////////////////////////////////////////////

func logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request %s %s %s %d\n",
			r.Method, r.RequestURI, r.RemoteAddr, requestIDCounter)
		next(w, r)
	}
}

func requestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// atomically increment the counter for the next request
		atomic.AddUint64(&requestIDCounter, 1)
		r = r.WithContext(context.WithValue(ctx, contextKeyRequestID, requestIDCounter))
		defer log.Printf("Finished handling http req. %d\n", requestIDCounter)

		r.Header.Set(HeaderRequestID, fmt.Sprintf("%d", requestIDCounter))

		next.ServeHTTP(w, r)
	}
}

// Post is middleware to reject requests that are _not_ POST method
func Post(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

/////////////////////////////////////////////////////////////////////
// Handlers
/////////////////////////////////////////////////////////////////////

func handler1(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	log.Printf("Got some data User: %d has a balance: %f\n", user.ID, user.Balance)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func handler2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := getRequestID(ctx)

	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			log.Println("Header:", name, value)
		}
	}

	u := User{ID: 123, Balance: 43.0, accountID: "ABCDEF123"}
	log.Printf("reqID: %d, User: %d, AccountID: %s", reqID, u.ID, u.accountID)

	json.NewEncoder(w).Encode(u)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

/////////////////////////////////////////////////////////////////////
// Server
/////////////////////////////////////////////////////////////////////

func main() {
	// initalise at 1 so 0 is a fail state
	atomic.AddUint64(&requestIDCounter, 1)

	mux := http.NewServeMux()

	mux.HandleFunc("/handler1", Post(requestID(logging(handler1))))
	mux.HandleFunc("/handler2", requestID(logging(handler2)))
	mux.HandleFunc("/health", healthCheckHandler)

	const address string = ":8000"
	var timeout time.Duration = 5 * time.Second

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	server := http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", address)
		serverErrors <- server.ListenAndServeTLS("localhost.crt", "localhost.key")
	}()

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		LogError("server", err)
		os.Exit(1)

	case sig := <-shutdown:
		log.Printf("main : %v : Start shutdown", sig)
		log.Println("main : shutting down fake DB connection gracefully...")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := server.Shutdown(ctx)
		if err != nil {
			LogError(
				fmt.Sprintf("main : Graceful shutdown did not complete in %v", timeout),
				err)
			err = server.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			log.Println("integrity issue caused shutdown")
			os.Exit(0)
		case err != nil:
			log.Fatal("could not stop server gracefully")
		}
	}
}
