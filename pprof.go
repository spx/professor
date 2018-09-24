package professor

import (
	"log"
	"net/http"
	"net/http/pprof"
)

var token = "securitytoken"

// init disables default handlers registered by importing net/http/pprof.
func init() {
	http.DefaultServeMux = http.NewServeMux()
}

func SetToken(t string) {
	token = t
}

// Handle adds standard pprof handlers to mux.
func Handle(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", checkToken(pprof.Index))
	mux.HandleFunc("/debug/pprof/cmdline", checkToken(pprof.Cmdline))
	mux.HandleFunc("/debug/pprof/profile", checkToken(pprof.Profile))
	mux.HandleFunc("/debug/pprof/symbol", checkToken(pprof.Symbol))
	mux.HandleFunc("/debug/pprof/trace", checkToken(pprof.Trace))
}

// NewServeMux builds a ServeMux and populates it with standard pprof handlers.
func NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	Handle(mux)
	return mux
}

// NewServer constructs a server at addr with the standard pprof handlers.
func NewServer(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: NewServeMux(),
	}
}

// ListenAndServe starts a server at addr with standard pprof handlers.
func ListenAndServe(addr string) error {
	return NewServer(addr).ListenAndServe()
}

// Launch a standard pprof server at addr.
func Launch(addr string) {
	go func() {
		log.Println(ListenAndServe(addr))
	}()
}

func checkToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("token") != token {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
}
