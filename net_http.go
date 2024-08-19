package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/justinas/alice"
)

func theClient() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// create a request instance
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, "https://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-My-Client", "Learning go")

	// make the request
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("unexpected status: got %v", res.StatusCode))
	}
	fmt.Println(res.Header.Get("Content-Type"))
	var data struct {
		UserID    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"compeleted"`
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", data)
}

type HelloHandler struct{}

func (hh HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Handler!\n"))
}

type application struct {
	logger *slog.Logger
}

// middleware
func (app application) RequestTimer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		dur := time.Since(start)
		app.logger.Info("request time", "path", r.URL.Path, "duration", dur, "ip", r.RemoteAddr)
	})
}

var securityMsg = []byte("You didn't give the secret password\n")

func TerribleSecurityProvider(password string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Secret-Password") != password {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(securityMsg)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func theServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello httprouter!\n"))
		})

	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		w.Write([]byte(fmt.Sprintf("Hello %s!\n", name)))
	})

	mux.HandleFunc("GET /time", func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now().Format(time.RFC3339)
		w.Write([]byte(currentTime))
	})

	// Since *http.ServeMux dispatches requests -> http.Handler
	// And *http.ServeMux implements http.Handler
	// We can create instance of *http.ServeMux with multiple related requests
	// and register it with a parent *http.ServeMux
	person := http.NewServeMux()
	person.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("greetings!\n"))
	})

	dog := http.NewServeMux()
	dog.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("good puppy!\n"))
	})

	mux.Handle("/person/", http.StripPrefix("/person", person))
	mux.Handle("/dog/", http.StripPrefix("/dog", dog))

	logOptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logHandler := slog.NewJSONHandler(os.Stderr, logOptions)
	logger := slog.New(logHandler)

	app := &application{
		logger: logger,
	}

	// terribleSecurity := TerribleSecurityProvider("GOPHER")

	// chain := alice.New(terribleSecurity, RequestTimer).ThenFunc(mux.ServeHTTP)
	chain := alice.New(app.RequestTimer).ThenFunc(mux.ServeHTTP)

	port := ":8000"
	s := http.Server{
		Addr:         port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      chain,
	}

	app.logger.Info("server started successfully", "port", port)
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}

func handler(rw http.ResponseWriter, req *http.Request) {
	rc := http.NewResponseController(rw)
	for i := 0; i < 10; i++ {
		_, err := rw.Write([]byte("asd"))
		if err != nil {
			slog.Error("error writing", "msg", err)
		}

		err = rc.Flush()
		if err != nil && !errors.Is(err, http.ErrNotSupported) {
			slog.Error("error flusing", "msg", err)
			return
		}
	}
}
