package admin

import "net/http"

func RegisterRoutes(h *Handler, mux *http.ServeMux) {
	// ========== HTMX Admin Routes ==========
	// Full admin page
	mux.HandleFunc("/admin", h.AdminPage)

	// Fragments for HTMX
	mux.HandleFunc("/admin/songs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.SongsFragment(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/admin/queue", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.QueueFragment(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/admin/queue/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.AddToQueueHTMX(w, r)
		} else if r.Method == http.MethodDelete {
			h.RemoveFromQueueHTMX(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/admin/songs/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			h.DeleteSongHTMX(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/admin/upload", h.UploadSongHTMX)
	mux.HandleFunc("/admin/play", h.PlayControl)
	mux.HandleFunc("/admin/now-playing", h.NowPlayingFragment)

	// ========== Legacy API Routes (JSON) ==========
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
