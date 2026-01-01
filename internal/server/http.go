package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andreasnainggolan/simple-gateway/internal/auth"
	"github.com/andreasnainggolan/simple-gateway/internal/proxy"
	"github.com/andreasnainggolan/simple-gateway/internal/router"
)

type Server struct {
	ListenAddr string
	Router     *router.Router
}

/*
New membuat instance HTTP server gateway.
*/
func New(listenAddr string, r *router.Router) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Router:     r,
	}
}

/*
Start menjalankan HTTP server dan handler utama gateway.
*/
func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		path := r.URL.Path
		method := r.Method

		log.Printf("[REQ] %s host=%s path=%s", method, host, path)

		// 1. Cari route yang cocok
		match, ok := s.Router.Match(host, path)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error":   "not_found",
				"message": "Route tidak ditemukan",
				"status":  404,
			})
			return
		}

		// 2. API KEY AUTH (jika diaktifkan)
		if match.Route.Protect != nil && match.Route.Protect.APIKey {
			if !auth.CheckAPIKey(r) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"error":   "unauthorized",
					"message": "API key tidak valid atau tidak ada",
					"status":  401,
				})
				return
			}
		}

		// 3. Reverse proxy ke backend
		p, err := proxy.New(match.Route.ForwardTo)
		if err != nil {
			log.Printf("[ERROR] invalid upstream: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		p.ServeHTTP(w, r)
	})

	log.Printf("[START] listening on %s", s.ListenAddr)
	return http.ListenAndServe(s.ListenAddr, mux)
}
