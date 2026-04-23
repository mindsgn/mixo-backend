package admin

import "net/http"

func RegisterRoutes(h *Handler, mux *http.ServeMux) {
	// Songs
	mux.HandleFunc("/api/songs", h.ListSongs)
	mux.HandleFunc("/api/songs/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.AddSong(w, r)
		} else if r.Method == http.MethodDelete {
			h.DeleteSong(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Upload
	mux.HandleFunc("/api/upload", h.UploadSong)

	// Queue
	mux.HandleFunc("/api/queue", h.GetQueue)
	mux.HandleFunc("/api/queue/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.AddToQueue(w, r)
		} else if r.Method == http.MethodDelete {
			h.RemoveFromQueue(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Now playing
	mux.HandleFunc("/api/now-playing", h.NowPlaying)

	// History
	mux.HandleFunc("/api/history", h.GetHistory)
}
